package pubsub

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type dummyConsumerGroupClaim struct {
	ch chan *sarama.ConsumerMessage
}

func (c *dummyConsumerGroupClaim) Topic() string                            { return "" }
func (c *dummyConsumerGroupClaim) Partition() int32                         { return 0 }
func (c *dummyConsumerGroupClaim) InitialOffset() int64                     { return 0 }
func (c *dummyConsumerGroupClaim) HighWaterMarkOffset() int64               { return 0 }
func (c *dummyConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage { return c.ch }

type dummyConsumerGroupSession struct{}

func (s *dummyConsumerGroupSession) Claims() map[string][]int32                               { return nil }
func (s *dummyConsumerGroupSession) MemberID() string                                         { return "" }
func (s *dummyConsumerGroupSession) GenerationID() int32                                      { return 0 }
func (s *dummyConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}
func (s *dummyConsumerGroupSession) Commit()                                                  {}
func (s *dummyConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}
func (s *dummyConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {}
func (s *dummyConsumerGroupSession) Context() context.Context                                 { return nil }

func TestConsumeClaimShouldExitWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	handler := consumerGroupHandler{ctx: ctx}

	consumeFinished := false
	go func() {
		_ = handler.ConsumeClaim(nil, &dummyConsumerGroupClaim{})
		consumeFinished = true
	}()

	cancel()
	time.Sleep(time.Second)
	require.True(t, consumeFinished)
}

func TestConsumeClaimShouldPassMessage(t *testing.T) {
	originalMessage := &sarama.ConsumerMessage{Topic: "test_topic"}
	var passedMessage *sarama.ConsumerMessage

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := consumerGroupHandler{
		ctx: ctx,
		onMessage: func(message *sarama.ConsumerMessage) {
			passedMessage = message
		},
	}

	ch := make(chan *sarama.ConsumerMessage)
	claim := &dummyConsumerGroupClaim{ch: ch}

	go func() {
		_ = handler.ConsumeClaim(&dummyConsumerGroupSession{}, claim)
	}()

	ch <- originalMessage
	time.Sleep(time.Second)
	require.NotNil(t, passedMessage)
	require.Equal(t, passedMessage.Topic, originalMessage.Topic)
}
