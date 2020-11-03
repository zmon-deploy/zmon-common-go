package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	retry := Retry(3)
	counter := 0
	_ = retry.Do(func() (bool, error) {
		counter++
		return false, nil
	})
	require.Equal(t, 3, counter)

	_ = retry.Do(func() (bool, error) {
		counter++
		if counter >= 5 {
			return true, nil
		}
		return false, nil
	})
	require.Equal(t, 5, counter)
}
