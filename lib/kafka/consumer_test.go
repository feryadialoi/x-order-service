package kafka

import (
	"context"
	"fmt"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/suite"
)

type ConsumerSuite struct {
	suite.Suite

	consumer *Consumer[customerDetail]
}

func (s *ConsumerSuite) SetupTest() {
	s.consumer = &Consumer[customerDetail]{
		consumerGroup: &testConsumerGroup{},
	}
}

func (s *ConsumerSuite) SetupSubTest() {
	s.consumer = &Consumer[customerDetail]{
		consumerGroup: &testConsumerGroup{},
	}
}

func (s *ConsumerSuite) TestStart() {
	s.Run("should start successfully", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			cancel()
		}()

		consumeClaim := func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
			return nil
		}

		s.consumer.handler = &testConsumerGroupHandler{
			consumeClaim: consumeClaim,
		}

		err := s.consumer.Start(ctx)
		s.NoError(err)
	})

	s.Run("should stop due to got error closed consumer group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			cancel()
		}()

		consumeClaim := func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
			return sarama.ErrClosedConsumerGroup
		}

		s.consumer.handler = &testConsumerGroupHandler{
			consumeClaim: consumeClaim,
		}

		err := s.consumer.Start(ctx)
		s.Error(err)
	})

	s.Run("should continue since not error closed consumer group", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			cancel()
		}()

		consumeClaim := func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
			return fmt.Errorf("test error")
		}

		s.consumer.handler = &testConsumerGroupHandler{
			consumeClaim: consumeClaim,
		}

		err := s.consumer.Start(ctx)
		s.NoError(err)
	})
}

func (s *ConsumerSuite) TestStop() {
	s.Run("should stop successfully", func() {
		s.consumer.consumerGroup = &testConsumerGroup{
			closeErr: nil,
		}

		err := s.consumer.Stop(context.Background())
		s.NoError(err)
	})

	s.Run("should log error when failed to close consumer group", func() {
		s.consumer.consumerGroup = &testConsumerGroup{
			closeErr: fmt.Errorf("test error"),
		}

		err := s.consumer.Stop(context.Background())
		s.Error(err)
	})
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}
