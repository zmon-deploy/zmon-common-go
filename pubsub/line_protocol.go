package pubsub

import (
	"bytes"
	protocol "github.com/influxdata/line-protocol"
	"github.com/pkg/errors"
	"time"
)

type metric struct {
	measurement string
	tags        map[string]string
	fields      map[string]interface{}
	tm          time.Time
}

func (m *metric) Time() time.Time {
	return m.tm
}

func (m *metric) Name() string {
	return m.measurement
}

func (m *metric) TagList() []*protocol.Tag {
	tagList := make([]*protocol.Tag, len(m.tags))
	counter := 0
	for key, value := range m.tags {
		tagList[counter] = &protocol.Tag{Key: key, Value: value}
		counter++
	}
	return tagList
}

func (m *metric) FieldList() []*protocol.Field {
	fieldList := make([]*protocol.Field, len(m.fields))
	counter := 0
	for key, value := range m.fields {
		fieldList[counter] = &protocol.Field{Key: key, Value: value}
		counter++
	}
	return fieldList
}

func encodeLineProtocol(measurement string, tags map[string]string, fields map[string]interface{}, tm time.Time) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := protocol.NewEncoder(buf)
	encoder.SetMaxLineBytes(1024)
	encoder.SetFieldTypeSupport(protocol.UintSupport)

	_, err := encoder.Encode(&metric{
		measurement: measurement,
		tags: tags,
		fields: fields,
		tm: tm,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode metric")
	}

	return buf.Bytes(), nil
}
