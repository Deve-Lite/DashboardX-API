package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Int struct {
	Int  int
	Null bool
	Set  bool
}

func NewInt(value int, null bool, set bool) Int {
	return Int{value, null, set}
}

func (t *Int) UnmarshalJSON(data []byte) error {
	t.Set = true

	if string(data) == "null" {
		t.Null = true
		return nil
	}

	var temp int
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	t.Int = temp
	t.Null = false
	return nil
}

func (t Int) Value() (driver.Value, error) {
	return int64(t.Int), nil
}

func (t *Int) Scan(v interface{}) error {
	if v == nil {
		*t = NewInt(0, true, true)
		return nil
	}

	if bv, err := driver.Int32.ConvertValue(v); err == nil {
		if v, ok := bv.(int); ok {
			*t = NewInt(v, false, true)
			return nil
		}
	}

	return errors.New("failed to scan Int")
}
