package stringutil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	require.Equal(t, "hello_world", ToSnakeCase("helloWorld"))
	require.Equal(t, "api_server", ToSnakeCase("APIServer"))
}

func TestStringContains(t *testing.T) {
	require.True(t, StringContains([]string{"a", "b", "c"}, "a"))
	require.False(t, StringContains([]string{"a", "b", "c"}, "d"))
}

func TestStrPtr(t *testing.T) {
	require.Equal(t, "abc", *StrPtr("abc"))
}

func TestHash(t *testing.T) {
	require.Equal(t, uint32(0x8bb5a57), Hash("asdf"))
}