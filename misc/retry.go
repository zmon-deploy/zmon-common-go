package misc

import "time"

type retryer struct {
	retryCount int
}

type retryFn func() bool

func (r *retryer) Do(callback retryFn) {
	r.DoWithDelay(callback, 0)
}

func (r *retryer) DoWithDelay(callback retryFn, delay time.Duration) {
	for i := 0; i < r.retryCount; i++ {
		success := callback()
		if success || i >= r.retryCount-1 {
			return
		} else {
			time.Sleep(delay)
		}
	}
}

func Retry(retryCount int) *retryer {
	return &retryer{retryCount}
}
