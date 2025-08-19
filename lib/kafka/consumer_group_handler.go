package kafka

import (
	"encoding/json"

	xlog "bitbucket.org/Amartha/go-x/log"
	"github.com/IBM/sarama"
)

type consumerGroupHandler[T any] struct {
	handler ConsumerHandler[T]
}

var _ sarama.ConsumerGroupHandler = new(consumerGroupHandler[any])

func (c *consumerGroupHandler[T]) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler[T]) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler[T]) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			if err := c.handle(sess, msg); err != nil {
				xlog.Error(sess.Context(), "failed to process message", xlog.Err(err))
			}
		case <-sess.Context().Done():
			xlog.Info(sess.Context(), "consumer session done")
			return nil
		}
	}
}

func (c *consumerGroupHandler[T]) unmarshal(value []byte) (T, error) {
	var val T

	// would always unmarshal using json
	// unless the value implements Unmarshaler
	// then it will be unmarshaled using the Unmarshal method
	switch v := any(val).(type) {
	case string:
		v = string(value)
	case Unmarshaler:
		if err := v.Unmarshal(value); err != nil {
			return val, err
		}
	default:
		if err := json.Unmarshal(value, &val); err != nil {
			return val, err
		}
	}

	return val, nil
}

func (c *consumerGroupHandler[T]) handle(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	defer sess.MarkMessage(msg, "")

	ctx := sess.Context()
	ctx = TopicToContext(ctx, msg.Topic)
	ctx = PartitionToContext(ctx, msg.Partition)
	ctx = OffsetToContext(ctx, msg.Offset)
	ctx = TimestampToContext(ctx, msg.Timestamp)

	value, err := c.unmarshal(msg.Value)
	if err != nil {
		return err
	}

	header := new(Header).FromSarama(msg.Headers)

	if err = c.handler.Handle(ctx, string(msg.Key), value, header); err != nil {
		return err
	}

	return nil
}
