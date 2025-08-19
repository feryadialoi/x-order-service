package kafka

import "github.com/IBM/sarama"

type testSyncProducer struct {
	sendMessage func(message *sarama.ProducerMessage) (partition int32, offset int64, err error)

	sendMessages struct {
		err error
	}

	close struct {
		err error
	}

	txnStatus sarama.ProducerTxnStatusFlag

	isTransactional struct {
		ok bool
	}

	beginTxn struct {
		err error
	}

	commitTxn struct {
		err error
	}

	abortTxn struct {
		err error
	}

	addOffsetsToTxn struct {
		err error
	}

	addMessageToTxn struct {
		err error
	}
}

var _ sarama.SyncProducer = new(testSyncProducer)

func (s *testSyncProducer) SendMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return s.sendMessage(message)
}

func (s *testSyncProducer) SendMessages(messages []*sarama.ProducerMessage) error {
	return s.sendMessages.err
}

func (s *testSyncProducer) Close() error {
	return s.close.err
}

func (s *testSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	return s.txnStatus
}

func (s *testSyncProducer) IsTransactional() bool {
	return s.isTransactional.ok
}

func (s *testSyncProducer) BeginTxn() error {
	return s.beginTxn.err
}

func (s *testSyncProducer) CommitTxn() error {
	return s.commitTxn.err
}

func (s *testSyncProducer) AbortTxn() error {
	return s.abortTxn.err
}

func (s *testSyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	return s.addOffsetsToTxn.err
}

func (s *testSyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	return s.addMessageToTxn.err
}
