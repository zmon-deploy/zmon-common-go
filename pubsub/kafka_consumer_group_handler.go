package pubsub

import (
	"context"
	"github.com/Shopify/sarama"
)

type consumerGroupHandler struct {
	ctx context.Context
	onMessage func(message *sarama.ConsumerMessage)
}

func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-h.ctx.Done():
			return nil // stop
		case message, ok := <-claim.Messages():
			if !ok {
				return nil // stop
			}
			h.onMessage(message)
			session.MarkMessage(message, "")
		}
	}
}
