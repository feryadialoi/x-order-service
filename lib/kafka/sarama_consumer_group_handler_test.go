package kafka

import "github.com/IBM/sarama"

type testConsumerGroupHandler struct {
	consumeClaim func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

var _ sarama.ConsumerGroupHandler = new(testConsumerGroupHandler)

func (c *testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return c.consumeClaim(session, claim)
}
