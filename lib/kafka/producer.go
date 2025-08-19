package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

type ProducerConfig struct {
	saramaConfig *sarama.Config
}

func defaultProducerConfig() ProducerConfig {
	config := sarama.NewConfig()
	config.Version = sarama.V3_4_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	return ProducerConfig{
		saramaConfig: config,
	}
}

type ProducerOption[T any] func(*ProducerConfig)

//go:generate mockgen -source=producer.go -package=mock -destination=mock/producer_mock.go

type Producer[T any] interface {
	Send(ctx context.Context, key string, value T, header Header, opts ...SendOption) error
}

type SyncProducer[T any] struct {
	producer sarama.SyncProducer
	topic    string
}

func NewSyncProducer[T any](addresses []string, topic string, opts ...ProducerOption[T]) (*SyncProducer[T], error) {
	cfg := defaultProducerConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	syncProducer, err := sarama.NewSyncProducer(addresses, cfg.saramaConfig)
	if err != nil {
		return nil, err
	}

	producer := SyncProducer[T]{
		producer: syncProducer,
		topic:    topic,
	}

	return &producer, nil
}

type SendConfig struct {
	Metadata  any
	Offset    int64
	Partition int32
}

type SendOption func(*SendConfig)

func WithSendMetadata(metadata any) SendOption {
	return func(cfg *SendConfig) {
		cfg.Metadata = metadata
	}
}

func WithSendOffset(offset int64) SendOption {
	return func(cfg *SendConfig) {
		cfg.Offset = offset
	}
}

func WithSendPartition(partition int32) SendOption {
	return func(cfg *SendConfig) {
		cfg.Partition = partition
	}
}

func (p *SyncProducer[T]) Send(ctx context.Context, key string, value T, header Header, opts ...SendOption) error {
	var sendConfig SendConfig

	for _, opt := range opts {
		opt(&sendConfig)
	}

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var headers []sarama.RecordHeader
	for k, v := range header {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	if _, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(val),
		Headers:   headers,
		Metadata:  sendConfig.Metadata,
		Offset:    sendConfig.Offset,
		Partition: sendConfig.Partition,
	}); err != nil {
		return err
	}

	return nil
}
