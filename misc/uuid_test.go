package misc

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestUUID(t *testing.T) {
	require.Len(t, UUIDString(), 36)
	require.Len(t, strings.Join(strings.Split(UUIDString(), "-"), ""), 32)
}
