package numberutil

import (
	"math/rand"
	"strconv"
)

func ConvertToFloat64(value interface{}) float64 {
	var result float64
	switch v := value.(type) {
	default:
		result = 0.0
	case uint:
		result = float64(v)
	case uint8:
		result = float64(v)
	case uint16:
		result = float64(v)
	case uint32:
		result = float64(v)
	case uint64:
		result = float64(v)
	case int:
		result = float64(v)
	case int8:
		result = float64(v)
	case int16:
		result = float64(v)
	case int32:
		result = float64(v)
	case int64:
		result = float64(v)
	case float32:
		result = float64(v)
	case string:
		n, _ := strconv.ParseFloat(v, 64)
		result = n
	case float64:
		result = v
	}
	return result
}

func GetRandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
