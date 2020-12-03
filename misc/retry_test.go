package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	retry := Retry(3)
	counter := 0
	retry.Do(func() bool {
		counter++
		return false
	})
	require.Equal(t, 3, counter)

	retry.Do(func() bool {
		counter++
		if counter >= 5 {
			return true
		}
		return false
	})
	require.Equal(t, 5, counter)
}
