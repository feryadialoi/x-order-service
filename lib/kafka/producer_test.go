package kafka

import (
	"context"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/suite"
)

type ProducerSuite struct {
	suite.Suite

	producer *SyncProducer[customerDetail]
}

func (s *ProducerSuite) SetupTest() {
	s.producer = &SyncProducer[customerDetail]{
		producer: &testSyncProducer{},
	}
}

func (s *ProducerSuite) SetupSubTest() {
	s.producer = &SyncProducer[customerDetail]{
		producer: &testSyncProducer{},
	}
}

func (s *ProducerSuite) TestSend() {
	s.Run("should send successfully", func() {
		s.producer.producer = &testSyncProducer{
			sendMessage: func(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
				key, err := message.Key.Encode()
				s.NoError(err)
				s.Equal("5000000001", string(key))

				value, err := message.Value.Encode()
				s.NoError(err)
				s.JSONEq(`{
					"customerNumber":"5000000001",
					"name":"customer",
					"sms":"081234567890"
				}`, string(value))

				s.Equal(1, len(message.Headers))
				s.Equal("key", string(message.Headers[0].Key))
				s.Equal("value", string(message.Headers[0].Value))

				s.Equal("metadata", message.Metadata)
				s.Equal(int32(0), message.Partition)
				s.Equal(int64(0), message.Offset)

				return 0, 0, nil
			},
		}

		err := s.producer.Send(context.Background(),
			"5000000001",
			customerDetail{
				CustomerNumber: "5000000001",
				Name:           "customer",
				SMS:            "081234567890",
			},
			Header{
				"key": "value",
			},
			WithSendMetadata("metadata"),
			WithSendPartition(0),
			WithSendOffset(0),
		)

		s.NoError(err)
	})
}

func TestProducerSuite(t *testing.T) {
	suite.Run(t, new(ProducerSuite))
}
