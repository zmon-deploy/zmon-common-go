package log

import (
	"container/list"
	"fmt"
	"github.com/cnpst/zmon-common-go/timeutil"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

type sentryEventHistory struct {
	Err      error
	LastSent time.Time
}

func newSentryEventHistory(err error, lastSent time.Time) *sentryEventHistory {
	return &sentryEventHistory{
		Err:      err,
		LastSent: lastSent,
	}
}

type SentryWrapper interface {
	CaptureException(error)
}

type sentryWrapper struct {
	sync.Mutex
	internalList *list.List
	reportFn     func(error)
	clock        timeutil.Clock
}

func newSentryWrapper(sentryDSN, ver string) (SentryWrapper, error) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:     sentryDSN,
		Release: fmt.Sprintf("zmon-consumer@%s", ver),
	}); err != nil {
		return nil, errors.Wrap(err, "failed to initialize sentry wrapper")
	}

	return &sentryWrapper{
		internalList: list.New(),
		reportFn: func(err error) {
			sentry.CaptureException(err)
			sentry.Flush(5 * time.Second)
		},
		clock: timeutil.NewClock(),
	}, nil
}

func (sw *sentryWrapper) CaptureException(err error) {
	sw.Lock()
	defer sw.Unlock()

	event := sw.getEvent(err)
	if event != nil {
		if sw.clock.Since(event.LastSent) > time.Hour {
			sw.removeEvent(err)
			sw.reportFn(err)
			sw.pushEvent(err)
		}
	} else {
		sw.reportFn(err)
		sw.pushEvent(err)
	}
}

func (sw *sentryWrapper) getEvent(err error) *sentryEventHistory {
	for e := sw.internalList.Front(); e != nil; e = e.Next() {
		evt := e.Value.(*sentryEventHistory)
		if evt.Err.Error() == err.Error() {
			return evt
		}
	}
	return nil
}

func (sw *sentryWrapper) removeEvent(err error) {
	var found *list.Element

	for e := sw.internalList.Front(); e != nil; e = e.Next() {
		evt := e.Value.(*sentryEventHistory)
		if evt.Err.Error() == err.Error() {
			found = e
			break
		}
	}

	sw.internalList.Remove(found)
}

func (sw *sentryWrapper) pushEvent(err error) {
	event := newSentryEventHistory(err, sw.clock.Now())
	sw.internalList.PushFront(event)
}
