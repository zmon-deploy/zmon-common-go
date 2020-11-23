package misc

import (
	"context"
	"time"
)

func Schedule(ctx context.Context, period time.Duration, fn func()) {
	go func() {
		ticker := time.NewTicker(period)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fn()
			}
		}
	}()
}

func ScheduleWithInitialDelay(ctx context.Context, period, initialDelay time.Duration, fn func()) {
	go func() {
		select {
		case <-time.After(initialDelay):
			Schedule(ctx, period, fn)
		case <-ctx.Done():
			return
		}
	}()
}
