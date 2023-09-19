package nullable

import (
	"encoding/json"
)

type StringsMap struct {
	Map  map[string]string
	Null bool
	Set  bool
}

func NewStringsMap(value map[string]string, null bool, set bool) StringsMap {
	return StringsMap{value, null, set}
}

func (t *StringsMap) UnmarshalJSON(data []byte) error {
	t.Set = true

	if string(data) == "null" {
		t.Null = true
		return nil
	}

	var temp map[string]string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	t.Map = temp
	t.Null = false
	return nil
}

// func (t StringsMap) Value() (driver.Value, error) {
// 	return int64(t.StringsMap), nil
// }

// func (t *StringsMap) Scan(v interface{}) error {
// 	if v == nil {
// 		*t = NewStringsMap(0, true, true)
// 		return nil
// 	}

// 	if bv, err := driver.StringsMap32.ConvertValue(v); err == nil {
// 		if v, ok := bv.(int); ok {
// 			*t = NewStringsMap(v, false, true)
// 			return nil
// 		}
// 	}

// 	return errors.New("failed to scan String")
// }
