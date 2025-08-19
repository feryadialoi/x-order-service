package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IBM/sarama"
)

type Offset int

const (
	OffsetOldest Offset = iota
	OffsetNewest
)

func (o Offset) Sarama() int64 {
	return map[Offset]int64{
		OffsetOldest: sarama.OffsetOldest,
		OffsetNewest: sarama.OffsetNewest,
	}[o]
}

type ConsumerConfig struct {
	saramaConfig *sarama.Config
}

func defaultConsumerConfig() ConsumerConfig {
	config := sarama.NewConfig()
	config.Version = sarama.V3_4_0_0

	return ConsumerConfig{
		saramaConfig: config,
	}
}

type ConsumerOption[T any] func(*ConsumerConfig)

func WithOffset[T any](offset Offset) ConsumerOption[T] {
	return func(cfg *ConsumerConfig) {
		if cfg.saramaConfig != nil {
			cfg.saramaConfig.Consumer.Offsets.Initial = offset.Sarama()
		}
	}
}

type Consumer[T any] struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	groupID       string
	handler       sarama.ConsumerGroupHandler
}

func NewConsumer[T any](addresses []string, topic string, groupID string, consumerHandler ConsumerHandler[T], opts ...ConsumerOption[T]) (*Consumer[T], error) {
	cfg := defaultConsumerConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	consumerGroup, err := sarama.NewConsumerGroup(addresses, groupID, cfg.saramaConfig)
	if err != nil {
		return nil, err
	}

	handler := consumerGroupHandler[T]{
		handler: consumerHandler,
	}

	consumer := Consumer[T]{
		consumerGroup: consumerGroup,
		topic:         topic,
		groupID:       groupID,
		handler:       &handler,
	}

	return &consumer, nil
}

func (c *Consumer[T]) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := c.consumerGroup.Consume(ctx, []string{c.topic}, c.handler)

			// consumer group is closed, it means the loop should stop
			if err != nil && errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return err
			}

			if err != nil {
				slog.Error("error consume", slog.String("error", err.Error()))
			}
		}
	}
}

func (c *Consumer[T]) Stop(ctx context.Context) error {
	err := c.consumerGroup.Close()
	if err != nil {
		slog.Error("error consume", slog.String("error", err.Error()))
	}

	return err
}
