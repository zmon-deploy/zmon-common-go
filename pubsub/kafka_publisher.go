package pubsub

import (
	"context"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaPublisher interface {
	Publish(ctx context.Context, messages ...kafka.Message) error
	PublishMetric(ctx context.Context, topic, measurement string, tags map[string]string, fields map[string]interface{}, tm time.Time) error
	Close() error
}

func NewKafkaPublisher(brokers []string, topic string) KafkaPublisher {
	return &kafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
			Async:    false,
		},
	}
}

type kafkaPublisher struct {
	writer *kafka.Writer
}

func (p *kafkaPublisher) Publish(ctx context.Context, messages ...kafka.Message) error {
	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return errors.Wrap(err, "failed to write messages")
	}

	return nil
}

func (p *kafkaPublisher) PublishMetric(ctx context.Context, topic, measurement string, tags map[string]string, fields map[string]interface{}, tm time.Time) error {
	encoded, err := encodeLineProtocol(measurement, tags, fields, tm)
	if err != nil {
		return errors.Wrap(err, "failed to encode line protocol")
	}

	message := kafka.Message{
		Topic: topic,
		Value: encoded,
	}

	return p.Publish(ctx, message)
}

func (p *kafkaPublisher) Close() error {
	if err := p.writer.Close(); err != nil {
		return errors.Wrap(err, "failed to close writer")
	}
	return nil
}
