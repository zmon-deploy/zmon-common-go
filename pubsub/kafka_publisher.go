package pubsub

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"time"
)

type KafkaPublisher interface {
	Publish(messages []*sarama.ProducerMessage) error
	PublishMetric(topic, measurement string, tags map[string]string, fields map[string]interface{}, tm time.Time) error
	Close() error
}

func NewKafkaPublisher(brokers []string, config *sarama.Config) (KafkaPublisher, error) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kafka sync producer")
	}

	return &kafkaPublisher{
		producer: producer,
	}, nil
}

func DefaultProducerConfig(clientID string) *sarama.Config {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.ClientID = clientID
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Compression = sarama.CompressionNone
	kafkaConfig.Producer.Retry.Max = 3
	kafkaConfig.Producer.Return.Successes = true

	return kafkaConfig
}

type kafkaPublisher struct {
	producer sarama.SyncProducer
}

func (p *kafkaPublisher) Publish(messages []*sarama.ProducerMessage) error {
	if err := p.producer.SendMessages(messages); err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (p *kafkaPublisher) PublishMetric(topic, measurement string, tags map[string]string, fields map[string]interface{}, tm time.Time) error {
	encoded, err := encodeLineProtocol(measurement, tags, fields, tm)
	if err != nil {
		return errors.Wrap(err, "failed to encode line protocol")
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(encoded),
	}

	return p.Publish([]*sarama.ProducerMessage{message})
}

func (p *kafkaPublisher) Close() error {
	return p.producer.Close()
}
