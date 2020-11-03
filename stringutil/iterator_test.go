package stringutil

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStringIterator_Default1(t *testing.T) {
	callCount := 0
	iter := NewStringIterator(func() ([]string, error) {
		callCount++
		if callCount == 1 {
			return []string{"a", "b"}, nil
		} else {
			return []string{}, nil
		}
	})
	require.Equal(t, true, iter.Next())
	require.Equal(t, "a", iter.Value())
	require.Equal(t, true, iter.Next())
	require.Equal(t, "b", iter.Value())
	require.Equal(t, false, iter.Next())
	require.NoError(t, iter.Err())
}

func TestStringIterator_Default2(t *testing.T) {
	callCount := 0
	iter := NewStringIterator(func() ([]string, error) {
		callCount++
		if callCount == 1 {
			return []string{"a", "b"}, nil
		} else if callCount == 2 {
			return []string{"c", "d"}, nil
		} else {
			return []string{}, nil
		}
	})
	require.Equal(t, true, iter.Next())
	require.Equal(t, "a", iter.Value())
	require.Equal(t, true, iter.Next())
	require.Equal(t, "b", iter.Value())
	require.Equal(t, true, iter.Next())
	require.Equal(t, "c", iter.Value())
	require.Equal(t, true, iter.Next())
	require.Equal(t, "d", iter.Value())
	require.Equal(t, false, iter.Next())
	require.NoError(t, iter.Err())
}

func TestStringIterator_Error1(t *testing.T) {
	iter := NewStringIterator(func() ([]string, error) {
		return nil, errors.New("test error")
	})
	require.Equal(t, false, iter.Next())
	require.Error(t, iter.Err())
}

func TestStringIterator_Error2(t *testing.T) {
	callCount := 0
	iter := NewStringIterator(func() ([]string, error) {
		callCount++
		if callCount == 1 {
			return []string{"a", "b"}, nil
		} else {
			return nil, errors.New("test error")
		}
	})
	require.Equal(t, true, iter.Next())
	require.Equal(t, "a", iter.Value())
	require.Equal(t, true, iter.Next())
	require.Equal(t, "b", iter.Value())
	require.Equal(t, false, iter.Next())
	require.Error(t, iter.Err())
}
