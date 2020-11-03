package stringutil

import (
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestStringSet(t *testing.T) {
	set := NewStringSet()
	require.Equal(t, []string{}, set.Values())

	set.Add("a", "b", "b")
	values := set.Values()
	sort.Strings(values)

	require.Equal(t, []string{"a", "b"}, values)
}
