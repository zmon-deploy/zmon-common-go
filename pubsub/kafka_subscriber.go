package pubsub

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/zmon-deploy/zmon-common-go/log"
	"github.com/pkg/errors"
	"time"
)

type OnKafkaMessage func(message *sarama.ConsumerMessage)

type KafkaSubscriber interface {
	Subscribe(topics []string, consumerGroup string, onMessage OnKafkaMessage) error
	Close()
}

func NewKafkaSubscriber(brokers []string, config *sarama.Config, logger log.Logger) KafkaSubscriber {
	return &kafkaSubscriber{
		logger:  log.NonNullLogger(logger),
		brokers: brokers,
		config:  config,
	}
}

func DefaultConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 10 * time.Second
	config.Consumer.Offsets.Retry.Max = 10

	return config
}

type kafkaSubscriber struct {
	logger  log.Logger
	brokers []string
	config  *sarama.Config
	cancel  context.CancelFunc
}

func (s *kafkaSubscriber) Subscribe(topics []string, consumerGroup string, onMessage OnKafkaMessage) error {
	var ctx context.Context
	ctx, s.cancel = context.WithCancel(context.Background())

	if err := s.consumeMessage(ctx, topics, consumerGroup, onMessage); err != nil {
		return err
	}

	return nil
}

func (s *kafkaSubscriber) consumeMessage(ctx context.Context, topics []string, consumerGroup string, onMessage OnKafkaMessage) error {
	handler := &consumerGroupHandler{
		ctx:       ctx,
		onMessage: onMessage,
	}

	consumer, err := sarama.NewConsumerGroup(s.brokers, consumerGroup, s.config)
	if err != nil {
		return errors.Wrap(err, "failed to create consumer")
	}

	go func() {
		for ctx.Err() == nil {
			if err := consumer.Consume(ctx, topics, handler); err != nil {
				if err == sarama.ErrUnknown {
					// ignore, because it's often just noise
				} else {
					s.logger.Errorf("failed to consume with client: %s", err.Error())
				}
				time.Sleep(time.Second)
			}
		}

		if err := consumer.Close(); err != nil {
			s.logger.Errorf("failed to close consumer client: %s", err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-consumer.Errors():
				s.logger.Errorf("error found on kafka consumer: %s", err.Error())
			}
		}
	}()

	return nil
}

func (s *kafkaSubscriber) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}
