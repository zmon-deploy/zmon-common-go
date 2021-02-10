package gormutil

import (
	"github.com/zmon-deploy/zmon-common-go/stringutil"
	"reflect"
	"strings"
)

func Columns(entity interface{}) []string {
	target := reflect.TypeOf(entity)

	var columns []string
	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		column := extractColumnFromField(field)
		columns = append(columns, column)
	}

	return columns
}

func extractColumnFromField(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")
	values := strings.Split(gormTag, ";")
	for _, value := range values {
		if strings.HasPrefix(value, "column:") {
			return value[len("column:"):]
		}
	}
	return stringutil.ToSnakeCase(field.Name)
}

