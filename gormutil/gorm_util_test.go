package gormutil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type testEntity struct {
	ID int64 `gorm:"column:id;primary_key"`
	Name string `gorm:"column:name"`
	CustomerEmail string `gorm:"column:customer_email"`
}

func TestColumns(t *testing.T) {
	require.Equal(t, []string{"id", "name", "customer_email"}, Columns(testEntity{}))
}
