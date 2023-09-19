package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Bool struct {
	Bool bool
	Null bool
	Set  bool
}

func NewBool(value bool, null bool, set bool) Bool {
	return Bool{value, null, set}
}

func (t *Bool) UnmarshalJSON(data []byte) error {
	t.Set = true

	if string(data) == "null" {
		t.Null = true
		return nil
	}

	var temp bool
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	t.Bool = temp
	t.Null = false
	return nil
}

func (t Bool) Value() (driver.Value, error) {
	return bool(t.Bool), nil
}

func (t *Bool) Scan(v interface{}) error {
	if v == nil {
		*t = NewBool(false, true, true)
		return nil
	}

	if bv, err := driver.Bool.ConvertValue(v); err == nil {
		if v, ok := bv.(bool); ok {
			*t = NewBool(v, false, true)
			return nil
		}
	}

	return errors.New("failed to scan String")
}
