package pubsub

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestEncodeLineProtocol(t *testing.T) {
	measurement := "test_measurement"
	tags := map[string]string{
		"host": "vm-01",
		"name": "hello",
	}
	fields := map[string]interface{}{
		"cpu":  0.123,
		"mem":  2222,
		"disk": 0.23,
		"log":  "helloWorld",
	}
	tm := time.Now()

	encoded, err := encodeLineProtocol(measurement, tags, fields, tm)
	require.NoError(t, err)

	result := string(encoded)
	splited := strings.Split(result, " ")
	require.Equal(t, 3, len(splited))

	measurementWithTags := strings.Split(splited[0], ",")
	actualMeasurement := measurementWithTags[0]
	actualTags := measurementWithTags[1:len(measurementWithTags)]
	sort.Strings(actualTags)

	require.Equal(t, measurement, actualMeasurement)
	require.Equal(t, "host=vm-01", actualTags[0])
	require.Equal(t, "name=hello", actualTags[1])

	actualFields := strings.Split(splited[1], ",")
	sort.Strings(actualFields)

	require.Equal(t, "cpu=0.123", actualFields[0])
	require.Equal(t, "disk=0.23", actualFields[1])
	require.Equal(t, "log=\"helloWorld\"", actualFields[2])
	require.Equal(t, "mem=2222i", actualFields[3])

	actualTime := splited[2]

	require.Equal(t, fmt.Sprintf("%v\n", tm.UnixNano()), actualTime)
}
