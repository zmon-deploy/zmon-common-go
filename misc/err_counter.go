package misc

import (
	"github.com/zmon-deploy/zmon-common-go/timeutil"
	"go.uber.org/multierr"
	"time"
)

// 일정 시간 (duration) 동안 일정 건수 (threshold) 이상의 err 발생 여부 카운팅
type ErrorCounter struct {
	entryRoot *entry
	count     int
	duration  time.Duration
	threshold int
	clock     timeutil.Clock
}

func NewErrorCounter(duration time.Duration, threshold int) *ErrorCounter {
	return &ErrorCounter{
		entryRoot: nil,
		count:     0,
		duration:  duration,
		threshold: threshold,
		clock:     timeutil.NewClock(),
	}
}

func (c *ErrorCounter) Put(err error) {
	now := c.clock.Now()
	newEntry := &entry{err: err, tm: now, next: nil}

	expireBaseLine := now.Add(-c.duration)
	ptr := c.entryRoot
	for {
		// if empty
		if ptr == nil {
			c.count++
			c.entryRoot = newEntry
			return
		}

		// drop if expired
		if ptr.tm.Before(expireBaseLine) {
			ptr = ptr.next
			c.entryRoot = ptr
			c.count--
			continue
		}

		// reached to end
		if ptr.next == nil {
			ptr.next = newEntry
			c.count++
			break
		}

		ptr = ptr.next
	}
}

func (c *ErrorCounter) IsOverThreshold() bool {
	return c.count >= c.threshold
}

func (c *ErrorCounter) Err() error {
	var err error
	ptr := c.entryRoot
	for {
		if ptr == nil {
			break
		}

		err = multierr.Append(err, ptr.err)
		ptr = ptr.next
	}
	return err
}

type entry struct {
	err  error
	tm   time.Time
	next *entry
}
