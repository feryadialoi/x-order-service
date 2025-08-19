package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type testConsumerGroupSession struct {
	claims       map[string][]int32
	memberID     string
	generationID int32
	markOffset   int32
	ctx          context.Context
}

var _ sarama.ConsumerGroupSession = new(testConsumerGroupSession)

func (c *testConsumerGroupSession) Claims() map[string][]int32 {
	return c.claims
}

func (c *testConsumerGroupSession) MemberID() string {
	return c.memberID
}

func (c *testConsumerGroupSession) GenerationID() int32 {
	return c.generationID
}

func (c *testConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}

func (c *testConsumerGroupSession) Commit() {
}

func (c *testConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

func (c *testConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
}

func (c *testConsumerGroupSession) Context() context.Context {
	return c.ctx
}
