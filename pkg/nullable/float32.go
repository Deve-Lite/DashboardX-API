package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Float32 struct {
	Float32 float32
	Null    bool
	Set     bool
}

func NewFloat32(value float32, null bool, set bool) Float32 {
	return Float32{value, null, set}
}

func (t *Float32) UnmarshalJSON(data []byte) error {
	t.Set = true

	if string(data) == "null" {
		t.Null = true
		return nil
	}

	var temp float32
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	t.Float32 = temp
	t.Null = false
	return nil
}

func (t Float32) Value() (driver.Value, error) {
	return int64(t.Float32), nil
}

func (t *Float32) Scan(v interface{}) error {
	if v == nil {
		*t = NewFloat32(0, true, true)
		return nil
	}

	if bv, err := driver.ValueConverter.ConvertValue(driver.DefaultParameterConverter, v); err == nil {
		if v, ok := bv.(float32); ok {
			*t = NewFloat32(v, false, true)
			return nil
		}
	}

	return errors.New("failed to scan String")
}
