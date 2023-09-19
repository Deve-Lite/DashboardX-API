package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type String struct {
	String string
	Null   bool
	Set    bool
}

func NewString(value string, null bool, set bool) String {
	return String{value, null, set}
}

func (t *String) UnmarshalJSON(data []byte) error {
	t.Set = true

	if string(data) == "null" {
		t.Null = true
		return nil
	}

	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	t.String = temp
	t.Null = false
	return nil
}

func (t String) Value() (driver.Value, error) {
	return string(t.String), nil
}

func (t *String) Scan(v interface{}) error {
	if v == nil {
		*t = NewString("", true, true)
		return nil
	}

	if bv, err := driver.String.ConvertValue(v); err == nil {
		if v, ok := bv.(string); ok {
			*t = NewString(v, false, true)
			return nil
		}
	}

	return errors.New("failed to scan String")
}
