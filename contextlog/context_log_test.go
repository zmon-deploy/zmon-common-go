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
