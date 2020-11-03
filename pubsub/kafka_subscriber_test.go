package pubsub

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHoldCancelFn(t *testing.T) {
	s := &kafkaSubscriber{
		cancels: map[uuid.UUID]context.CancelFunc{},
	}

	_, cancel := context.WithCancel(context.Background())
	clear := s.holdCancelFn(cancel)
	require.Equal(t, 1, len(s.cancels))
	clear()
	require.Equal(t, 0, len(s.cancels))
}
