package nullable

import (
	"database/sql"
	"encoding/json"
	"time"
)

// 기존 null 타입들에 json marshaling 관련 코드를 추가한 구조체

type NullInt32 struct {
	sql.NullInt32
}

func NewNullInt32(v *int32) NullInt32 {
	if v == nil {
		return NullInt32{sql.NullInt32{Valid: false}}
	}
	return NullInt32{sql.NullInt32{
		Int32: *v,
		Valid: true,
	}}
}

func NewNullInt32V(v int32) NullInt32 {
	return NewNullInt32(&v)
}

func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullInt32) UnmarshalJSON(data []byte) error {
	var x *int32
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int32 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type NullInt64 struct {
	sql.NullInt64
}

func NewNullInt64(v *int64) NullInt64 {
	if v == nil {
		return NullInt64{sql.NullInt64{Valid: false}}
	}
	return NullInt64{sql.NullInt64{
		Int64: *v,
		Valid: true,
	}}
}

func NewNullInt64V(v int64) NullInt64 {
	return NewNullInt64(&v)
}

func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type NullFloat64 struct {
	sql.NullFloat64
}

func NewNullFloat64(v *float64) NullFloat64 {
	if v == nil {
		return NullFloat64{sql.NullFloat64{Valid: false}}
	}
	return NullFloat64{
		sql.NullFloat64{
			Float64: *v,
			Valid:   true,
		},
	}
}

func NewNullFloat64V(v float64) NullFloat64 {
	return NewNullFloat64(&v)
}

func (v *NullFloat64) UnmarshalJSON(data []byte) error {
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Float64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type NullString struct {
	sql.NullString
}

func NewNullString(v *string) NullString {
	if v == nil {
		return NullString{sql.NullString{Valid: false}}
	}
	return NullString{sql.NullString{
		String: *v,
		Valid:  true,
	}}
}

func NewNullStringV(v string) NullString {
	return NewNullString(&v)
}

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}

type NullTime struct {
	sql.NullTime
}

func NewNullTime(t *time.Time) NullTime {
	if t == nil {
		return NullTime{sql.NullTime{Valid: false}}
	}
	return NullTime{sql.NullTime{
		Time:  *t,
		Valid: true,
	}}
}

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Time = *x
	} else {
		v.Valid = false
	}
	return nil
}