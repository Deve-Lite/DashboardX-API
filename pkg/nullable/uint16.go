package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Uint16 struct {
	Uint16 uint16
	Null   bool
	Set    bool
}

func NewUint16(value uint16, null bool, set bool) Uint16 {
	return Uint16{value, null, set}
}

func (t *Uint16) UnmarshalJSON(data []byte) error {
	t.Set = true

	if string(data) == "null" {
		t.Null = true
		return nil
	}

	var temp uint16
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	t.Uint16 = temp
	t.Null = false
	return nil
}

func (t Uint16) Value() (driver.Value, error) {
	return int64(t.Uint16), nil
}

func (t *Uint16) Scan(v interface{}) error {
	if v == nil {
		*t = NewUint16(0, true, true)
		return nil
	}

	if bv, err := driver.Int32.ConvertValue(v); err == nil {
		if v, ok := bv.(uint16); ok {
			*t = NewUint16(v, false, true)
			return nil
		}
	}

	return errors.New("failed to scan String")
}
