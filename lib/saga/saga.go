package saga

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/lib/kafka"
)

type Option[T any] struct {
	db      *sql.DB
	table   string
	broker  string
	topic   string
	groupID string
	machine StateMachine[T]
}

type OptionFunc[T any] func(option *Option[T])

func WithDatabase[T any](db *sql.DB, table string) OptionFunc[T] {
	return func(option *Option[T]) {
		option.db = db
		option.table = table
	}
}

func WithKafka[T any](broker string, topic string, groupID string) OptionFunc[T] {
	return func(option *Option[T]) {
		option.broker = broker
		option.topic = topic
		option.groupID = groupID
	}
}

func WithStateMachine[T any](machine StateMachine[T]) OptionFunc[T] {
	return func(option *Option[T]) {
		option.machine = machine
	}
}

type consumerHandler[T any] struct {
	db       *sql.DB
	table    string
	producer kafka.Producer[any]
	machine  StateMachine[T]
}

func (h *consumerHandler[T]) Handle(ctx context.Context, key string, value map[string]any, header kafka.Header) error {
	slog.Info("Handling message",
		slog.Any("key", key),
		slog.Any("value", value),
		slog.Any("header", header),
	)

	// Get the event type from the header
	eventTypeValue, ok := header["event_type"]
	if !ok {
		return nil // Skip messages without event type
	}

	// Convert the event type to a saga Event
	event := Event(eventTypeValue)

	// Get the event configurer from the state machine
	eventConfigurer := h.machine.EventConfigurer()

	// Check if we have a handler for this event
	eventHandler, ok := eventConfigurer[event]
	if !ok {
		return nil // Skip events we don't handle
	}

	// Marshal the value to JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Create a saga context for the generic type E
	c := Context[T]{
		ctx:     ctx,
		Message: valueBytes,
	}

	// Process the event to set correlation ID
	eventHandler(&c)

	// If no correlation ID was set, we can't proceed
	if c.Entity.CorrelationID == "" {
		slog.Warn("No correlation ID set",
			slog.Any("entity", c.Entity),
			slog.Any("event", event),
			slog.Any("value", value),
			slog.Any("header", header),
		)
		return nil
	}

	// Get the current state from the database or use initially state
	currentState := StateInitially

	// Try to get the current state from the database

	// First try to get just the state to determine if the record exists
	var stateStr string
	err = h.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`SELECT current_state FROM %s WHERE correlation_id = $1`, h.table),
		c.Entity.CorrelationID,
	).Scan(&stateStr)

	switch {
	case err == nil:
		currentState = State(stateStr)
		var entity Entity[T]
		err = h.db.QueryRowContext(
			ctx,
			`SELECT 
				correlation_id, 
				current_state,
				data,
				created_at,
				updated_at
			FROM `+h.table+` 
			WHERE correlation_id = $1`,
			c.Entity.CorrelationID,
		).Scan(
			&entity.CorrelationID,
			&entity.CurrentState,
			&entity.Data,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err == nil {
			c.Entity = entity
		}
	case !errors.Is(err, sql.ErrNoRows):
		return err
	}
	// If errors.Is(err, sql.ErrNoRows), we'll use the initial state

	// Get the event activities from the state machine
	eventActivities := h.machine.EventActivities()

	// Check if we have a transition for this state and event
	activity, ok := eventActivities[currentState]
	if !ok {
		return nil
	}

	transition, ok := activity[event]
	if !ok {
		return nil
	}

	// Execute the "Then" function if provided
	if transition.Then != nil {
		if err := transition.Then(&c); err != nil {
			return err
		}
	}

	// Publish a message if needed
	if transition.Publish != nil {
		slog.Info("Publishing message",
			slog.Any("transition", transition),
		)
		message := transition.Publish(&c)
		if message != nil {
			// Get the event type from the message
			if evt, ok := message.(contract.Event); ok {
				// Send the message
				if err := h.producer.Send(ctx, c.Entity.CorrelationID, message, kafka.Header{
					"event_type": evt.GetEventType(),
				}); err != nil {
					return err
				}
			}
		}
	}

	// Save the new state to the database

	// Check if we need to save or update the state
	var exists bool
	err = h.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM "+h.table+" WHERE correlation_id = $1)",
		c.Entity.CorrelationID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	// Set the new state
	c.Entity.CurrentState = string(transition.TransitionTo)

	slog.Info("Saving state",
		slog.Any("order_state", c.Entity),
	)

	now := time.Now().UTC()

	if exists {
		// Update existing record
		c.Entity.UpdatedAt = now

		_, err = h.db.ExecContext(
			ctx,
			`UPDATE `+h.table+`
			SET current_state = $1, 
				updated_at = $2
			WHERE correlation_id = $3`,
			c.Entity.CurrentState,
			c.Entity.UpdatedAt,
			c.Entity.CorrelationID,
		)
		if err != nil {
			slog.Error("Error updating record",
				slog.Any("entity", c.Entity),
				slog.Any("error", err),
			)
			return err
		}
	} else {
		// Insert a new record
		c.Entity.CreatedAt = now
		c.Entity.UpdatedAt = now

		_, err = h.db.ExecContext(
			ctx,
			`INSERT INTO `+h.table+` (
				correlation_id, 
				current_state, 
				data, 
				created_at, 
				updated_at
			) VALUES ($1, $2, $3, $4, $5)`,
			c.Entity.CorrelationID,
			c.Entity.CurrentState,
			c.Entity.Data,
			c.Entity.CreatedAt,
			c.Entity.UpdatedAt,
		)
		if err != nil {
			slog.Error("Error inserting record",
				slog.Any("entity", c.Entity),
				slog.Any("error", err),
			)
			return err
		}
	}

	return nil
}

type Saga struct {
	consumer *kafka.Consumer[map[string]any]
}

func NewSaga[T any](options ...OptionFunc[T]) *Saga {
	var option Option[T]
	for _, opt := range options {
		opt(&option)
	}

	producer, err := kafka.NewSyncProducer[any]([]string{option.broker}, option.topic)
	if err != nil {
		panic(err)
	}

	consumer, err := kafka.NewConsumer[map[string]any]([]string{option.broker}, option.topic, option.groupID, &consumerHandler[T]{
		db:       option.db,
		table:    option.table,
		producer: producer,
		machine:  option.machine,
	}, kafka.WithOffset[map[string]any](kafka.OffsetOldest))
	if err != nil {
		panic(err)
	}

	return &Saga{
		consumer: consumer,
	}
}

func (s *Saga) Start(ctx context.Context) error {
	return s.consumer.Start(ctx)
}

func (s *Saga) Stop(ctx context.Context) error {
	return s.consumer.Stop(ctx)
}
