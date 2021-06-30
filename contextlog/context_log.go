package contextlog

import (
	"context"
	"fmt"
	"strings"
)

const contextFieldsKey = "contextFields"

func WithContextField(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextFieldsKey, newContextFields())
}

func RestoreFromContextFields(fields *ContextFields) context.Context {
	ctx := WithContextField(context.Background())

	if fields != nil {
		fields.ForEach(func(key string, value interface{}) {
			AppendField(ctx, key, value)
		})
	}

	return ctx
}

func FieldsFromContext(ctx context.Context) *ContextFields {
	value := ctx.Value(contextFieldsKey)
	if value == nil {
		return nil
	}

	return value.(*ContextFields)
}

func FieldsFromContextAsLine(ctx context.Context) string {
	out := strings.Builder{}
	cnt := 0

	fields := FieldsFromContext(ctx)
	fields.ForEach(func(key string, value interface{}) {
		out.WriteString(fmt.Sprintf("[%s: %v]", key, value))
		cnt++
		if cnt < fields.Len() {
			out.WriteString(", ")
		}
	})

	return out.String()
}

func AppendField(ctx context.Context, key string, value interface{}) {
	f := FieldsFromContext(ctx)
	if f != nil {
		f.Append(key, value)
	}
}

func AppendFields(ctx context.Context, keys []string, values ...interface{}) {
	for i, value := range values {
		AppendField(ctx, keys[i], value)
	}
}

type ContextFields struct {
	Data        map[string]interface{} `json:"data"`
	KeysInOrder []string               `json:"keysInOrder"`
}

func newContextFields() *ContextFields {
	return &ContextFields{
		Data:        map[string]interface{}{},
		KeysInOrder: []string{},
	}
}

func (f *ContextFields) Append(key string, value interface{}) {
	f.KeysInOrder = append(f.KeysInOrder, key)
	f.Data[key] = value
}

func (f *ContextFields) ForEach(fn func(key string, value interface{})) {
	for _, key := range f.KeysInOrder {
		fn(key, f.Data[key])
	}
}

func (f *ContextFields) Len() int {
	return len(f.KeysInOrder)
}

