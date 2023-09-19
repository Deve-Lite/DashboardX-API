package nullable

import (
	"encoding/json"
)

type Nullable[T any] struct {
	Value T
	Null  bool
	Set   bool
}

func NewNullable[T any](value T, null bool, set bool) Nullable[T] {
	return Nullable[T]{value, null, set}
}

func (i *Nullable[any]) UnmarshalJSON(data []byte) error {
	i.Set = true

	if string(data) == "null" {
		i.Null = true
		return nil
	}

	var temp any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Null = false
	return nil
}
