package timeutil

import "time"

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type clock struct{}

func NewClock() Clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}

func (c *clock) Since(tm time.Time) time.Duration {
	return time.Since(tm)
}

type fixedClock struct {
	now time.Time
}

func NewFixedClock(now time.Time) Clock {
	return &fixedClock{
		now: now,
	}
}

func (c *fixedClock) Now() time.Time {
	return c.now
}

func (c *fixedClock) Since(tm time.Time) time.Duration {
	return c.now.Sub(tm)
}
