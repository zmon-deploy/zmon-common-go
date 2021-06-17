package pubsub

import (
	"context"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"time"
)

type OnKafkaMessage func(message kafka.Message)

type KafkaSubscriber interface {
	Subscribe(ctx context.Context, onMessage OnKafkaMessage) error
	Close() error
}

func NewKafkaSubscriber(config kafka.ReaderConfig) KafkaSubscriber {
	return &kafkaSubscriber{
		reader: kafka.NewReader(config),
	}
}

func DefaultConsumerConfig(brokers []string, consumerGroup string, topics []string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        consumerGroup,
		GroupTopics:    topics,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 5 * time.Second,
		GroupBalancers: []kafka.GroupBalancer{Balancer{}},
		StartOffset:    kafka.FirstOffset,
	}
}

type kafkaSubscriber struct {
	reader *kafka.Reader
}

func (s *kafkaSubscriber) Subscribe(ctx context.Context, onMessage OnKafkaMessage) error {
	for {
		m, err := s.reader.ReadMessage(ctx)
		if err != nil {
			if err.Error() != "context canceled" {
				return errors.Wrap(err, "failed to read message")
			}
			break
		}

		onMessage(m)
	}
	return nil
}

func (s *kafkaSubscriber) Close() error {
	if err := s.reader.Close(); err != nil {
		return errors.Wrap(err, "failed to close reader")
	}
	return nil
}
