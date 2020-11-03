package pubsub

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cnpst/zmon-common-go/log"
	"github.com/cnpst/zmon-common-go/misc"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type KafkaSubscriber interface {
	Subscribe(topics []string, consumerGroup string) (<-chan *sarama.ConsumerMessage, context.CancelFunc, error)
	Close()
}

func NewKafkaSubscriber(brokers []string, config *sarama.Config, logger log.Logger) KafkaSubscriber {
	return &kafkaSubscriber{
		logger:  log.NonNullLogger(logger),
		brokers: brokers,
		config:  config,
		cancels: map[uuid.UUID]context.CancelFunc{},
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
	cancels map[uuid.UUID]context.CancelFunc
}

func (s *kafkaSubscriber) Subscribe(topics []string, consumerGroup string) (<-chan *sarama.ConsumerMessage, context.CancelFunc, error) {
	output := make(chan *sarama.ConsumerMessage)
	ctx, cancel := context.WithCancel(context.Background())

	if err := s.consumeMessage(ctx, topics, consumerGroup, output); err != nil {
		return nil, nil, err
	}

	return output, s.holdCancelFn(cancel), nil
}

func (s *kafkaSubscriber) consumeMessage(ctx context.Context, topics []string, consumerGroup string, output chan<- *sarama.ConsumerMessage) error {
	handler := &consumerGroupHandler{
		ctx: ctx,
		onMessage: func(message *sarama.ConsumerMessage) {
			output <- message
		},
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

func (s *kafkaSubscriber) holdCancelFn(cancel context.CancelFunc) context.CancelFunc {
	uid := misc.UUID()
	s.cancels[uid] = cancel
	clear := func() {
		cancel()
		delete(s.cancels, uid)
	}
	return clear
}

func (s *kafkaSubscriber) Close() {
	cancels := s.cancels
	s.cancels = map[uuid.UUID]context.CancelFunc{}

	for _, cancel := range cancels {
		cancel()
	}
}
