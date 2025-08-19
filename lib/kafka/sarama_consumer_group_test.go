package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type testConsumerGroup struct {
	closeErr error
}

var _ sarama.ConsumerGroup = new(testConsumerGroup)

func (c *testConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	session := testConsumerGroupSession{}
	claim := testConsumerGroupClaim{}

	return handler.ConsumeClaim(&session, &claim)
}

func (c *testConsumerGroup) Errors() <-chan error {
	return nil
}

func (c *testConsumerGroup) Close() error {
	return c.closeErr
}

func (c *testConsumerGroup) Pause(partitions map[string][]int32) {
}

func (c *testConsumerGroup) Resume(partitions map[string][]int32) {
}

func (c *testConsumerGroup) PauseAll() {
}

func (c *testConsumerGroup) ResumeAll() {
}
