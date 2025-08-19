package kafka

import (
	"context"
	"fmt"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/suite"
)

type customerDetail struct {
	CustomerNumber string `json:"customerNumber"`
	Name           string `json:"name"`
	SMS            string `json:"sms"`
}

type consumerHandlerFunc[T any] func(ctx context.Context, key string, value T, header Header) error

func (f consumerHandlerFunc[T]) Handle(ctx context.Context, key string, value T, header Header) error {
	return f(ctx, key, value, header)
}

type consumerHandlerSuite struct {
	suite.Suite

	handler *consumerGroupHandler[customerDetail]
}

func (s *consumerHandlerSuite) SetupTest() {
	handler := consumerHandlerFunc[customerDetail](func(ctx context.Context, key string, value customerDetail, header Header) error {
		return nil
	})

	s.handler = &consumerGroupHandler[customerDetail]{
		handler: handler,
	}
}

func (s *consumerHandlerSuite) SetupSubTest() {
	handler := consumerHandlerFunc[customerDetail](func(ctx context.Context, key string, value customerDetail, header Header) error {
		return nil
	})

	s.handler = &consumerGroupHandler[customerDetail]{
		handler: handler,
	}
}

func (s *consumerHandlerSuite) TestSetup() {
	s.Run("always do nothing", func() {
		sess := testConsumerGroupSession{}
		err := s.handler.Setup(&sess)
		s.NoError(err)
	})
}

func (s *consumerHandlerSuite) TestCleanup() {
	s.Run("always do nothing", func() {
		sess := testConsumerGroupSession{}
		err := s.handler.Cleanup(&sess)
		s.NoError(err)
	})
}

func (s *consumerHandlerSuite) TestConsumeClaim() {
	s.Run("processing message and stopping due to context done", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		session := testConsumerGroupSession{
			ctx: ctx,
		}

		messages := make(chan *sarama.ConsumerMessage)
		defer close(messages)

		claim := testConsumerGroupClaim{
			messages: messages,
		}

		go func() {
			messages <- &sarama.ConsumerMessage{
				Topic:     "queue.customer",
				Partition: 0,
				Offset:    1,
				Key:       []byte("5000000001"),
				Value: []byte(`{
					"customerNumber": "5000000001",
					"name": "customer 5000000001",
					"sms": "+6281234567890"
				}`),
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("type"),
						Value: []byte("customer_updated"),
					},
				},
			}
			cancel()
		}()

		s.handler.handler = consumerHandlerFunc[customerDetail](func(ctx context.Context, key string, value customerDetail, header Header) error {
			s.Equal("queue.customer", TopicFromContext(ctx))
			s.Equal(int32(0), PartitionFromContext(ctx))
			s.Equal(int64(1), OffsetFromContext(ctx))

			s.Equal("customer_updated", header.Get("type"))
			s.Equal("5000000001", value.CustomerNumber)
			s.Equal("customer 5000000001", value.Name)
			s.Equal("+6281234567890", value.SMS)

			return nil
		})

		err := s.handler.ConsumeClaim(&session, &claim)

		s.NoError(err)
	})

	s.Run("processing message got error and stopping due to context done", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		session := testConsumerGroupSession{
			ctx: ctx,
		}

		messages := make(chan *sarama.ConsumerMessage)
		defer close(messages)

		claim := testConsumerGroupClaim{
			messages: messages,
		}

		go func() {
			messages <- &sarama.ConsumerMessage{
				Topic:     "queue.customer",
				Partition: 0,
				Offset:    1,
				Key:       []byte("5000000001"),
				Value: []byte(`{
					"customerNumber": "5000000001",
					"name": "customer 5000000001",
					"sms": "+6281234567890"
				}`),
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("type"),
						Value: []byte("customer_updated"),
					},
				},
			}
			cancel()
		}()

		s.handler.handler = consumerHandlerFunc[customerDetail](func(ctx context.Context, key string, value customerDetail, header Header) error {
			s.Equal("queue.customer", TopicFromContext(ctx))
			s.Equal(int32(0), PartitionFromContext(ctx))
			s.Equal(int64(1), OffsetFromContext(ctx))

			s.Equal("customer_updated", header.Get("type"))
			s.Equal("5000000001", value.CustomerNumber)
			s.Equal("customer 5000000001", value.Name)
			s.Equal("+6281234567890", value.SMS)

			// some business logic

			return fmt.Errorf("failed to process message")
		})

		err := s.handler.ConsumeClaim(&session, &claim)

		s.NoError(err)
	})
}

func (s *consumerHandlerSuite) Test_handle() {}

func Test_consumerHandlerSuite(t *testing.T) {
	suite.Run(t, new(consumerHandlerSuite))
}
