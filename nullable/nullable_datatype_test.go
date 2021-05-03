package nullable

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNullInt32Json(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "int32_not_null",
			input:    struct{ Value NullInt32 }{NewNullInt32V(123)},
			expected: "{\"Value\":123}",
		},
		{
			name:     "int32_null",
			input:    struct{ Value NullInt32 }{NewNullInt32(nil)},
			expected: "{\"Value\":null}",
		},
		{
			name:     "int64_not_null",
			input:    struct{ Value NullInt64 }{NewNullInt64V(123)},
			expected: "{\"Value\":123}",
		},
		{
			name:     "int64_null",
			input:    struct{ Value NullInt64 }{NewNullInt64(nil)},
			expected: "{\"Value\":null}",
		},
		{
			name:     "float64_not_null",
			input:    struct{ Value NullFloat64 }{NewNullFloat64V(123.45)},
			expected: "{\"Value\":123.45}",
		},
		{
			name:     "float64_null",
			input:    struct{ Value NullFloat64 }{NewNullFloat64(nil)},
			expected: "{\"Value\":null}",
		},
		{
			name:     "string_not_null",
			input:    struct{ Value NullString }{NewNullStringV("asdf")},
			expected: "{\"Value\":\"asdf\"}",
		},
		{
			name:     "string_null",
			input:    struct{ Value NullString }{NewNullString(nil)},
			expected: "{\"Value\":null}",
		},
		{
			name:     "time_not_ull",
			input:    struct{ Value NullTime }{NewNullTime(&time.Time{})},
			expected: "{\"Value\":\"0001-01-01T00:00:00Z\"}",
		},
		{
			name:     "time_null",
			input:    struct{ Value NullTime }{NewNullTime(nil)},
			expected: "{\"Value\":null}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := json.Marshal(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(bytes))
		})
	}
}

func TestAsPtr(t *testing.T) {
	now := time.Now()

	tests := []struct {
		input interface{}
		isNil bool
	}{
		{input: NewNullInt32(nil).AsPtr(), isNil: true},
		{input: NewNullInt32V(123).AsPtr(), isNil: false},
		{input: NewNullInt64(nil).AsPtr(), isNil: true},
		{input: NewNullInt64V(123).AsPtr(), isNil: false},
		{input: NewNullFloat64(nil).AsPtr(), isNil: true},
		{input: NewNullFloat64V(123).AsPtr(), isNil: false},
		{input: NewNullString(nil).AsPtr(), isNil: true},
		{input: NewNullStringV("asdf").AsPtr(), isNil: false},
		{input: NewNullTime(nil).AsPtr(), isNil: true},
		{input: NewNullTime(&now).AsPtr(), isNil: false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i + 1), func(t *testing.T) {
			if tc.isNil {
				require.Nil(t, tc.input)
			} else {
				require.NotNil(t, tc.input)
			}
		})
	}
}
