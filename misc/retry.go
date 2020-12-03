package misc

type retryer struct {
	retryCount int
}

type retryFn func() bool

func (r *retryer) Do(callback retryFn) {
	for i := 0; i < r.retryCount; i++ {
		success := callback()
		if success || i >= r.retryCount-1 {
			return
		}
	}
}

func Retry(retryCount int) *retryer {
	return &retryer{retryCount}
}
