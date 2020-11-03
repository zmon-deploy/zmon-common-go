package misc

type RetryFailedError struct{}

func (e RetryFailedError) Error() string {
	return "retry failed"
}

type retryer struct {
	retryCount int
}

type retryFn func() (bool, error)

func (r *retryer) Do(callback retryFn) error {
	for i := 0; i < r.retryCount; i++ {
		success, err := callback()

		if err != nil && i == r.retryCount-1 {
			return err
		}
		if success == true {
			return nil
		}
	}

	return RetryFailedError{}
}

func Retry(retryCount int) *retryer {
	return &retryer{retryCount}
}
