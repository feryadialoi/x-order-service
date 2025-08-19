package kafka

import (
	"github.com/IBM/sarama"
)

type testConsumerGroupClaim struct {
	topic               string
	partition           int32
	initialOffset       int64
	highWaterMarkOffset int64
	messages            chan *sarama.ConsumerMessage
}

var _ sarama.ConsumerGroupClaim = new(testConsumerGroupClaim)

func (c *testConsumerGroupClaim) Topic() string {
	return c.topic
}

func (c *testConsumerGroupClaim) Partition() int32 {
	return c.partition
}

func (c *testConsumerGroupClaim) InitialOffset() int64 {
	return c.initialOffset
}

func (c *testConsumerGroupClaim) HighWaterMarkOffset() int64 {
	return c.highWaterMarkOffset
}

func (c *testConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	return c.messages
}
