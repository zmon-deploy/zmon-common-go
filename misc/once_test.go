package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpirableOnce(t *testing.T) {
	once := ExpirableOnce{}
	counter := 0

	once.Do(func() { counter++ })
	once.Do(func() { counter++ })

	require.Equal(t, 1, counter)

	once.Expire()
	once.Do(func() { counter++ })

	require.Equal(t, 2, counter)
}
