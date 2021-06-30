package contextlog

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithLog(t *testing.T) {
	ctx := WithContextField(context.Background())
	AppendField(ctx, "first log", "hello")

	func(inCtx context.Context) {
		AppendField(inCtx, "second log", "world")
	}(ctx)

	AppendFields(ctx, []string{"another", "log"}, true, false)

	fields := FieldsFromContextAsLine(ctx)
	require.Equal(t, "[first log: hello], [second log: world], [another: true], [log: false]", fields)
}

func TestClone(t *testing.T) {
	originalCtx := WithContextField(context.Background())
	AppendField(originalCtx, "1", "1")
	AppendField(originalCtx, "who", "original")
	AppendField(originalCtx, "2", "2")

	clonedCtx := Clone(originalCtx)
	AppendField(clonedCtx, "who", "cloned")

	require.Equal(t, "[1: 1], [who: original], [2: 2]", FieldsFromContextAsLine(originalCtx))
	require.Equal(t, "[1: 1], [2: 2], [who: cloned]", FieldsFromContextAsLine(clonedCtx))
}
