package proto

import (
	"bytes"
	"encoding/json"
)

// structUnmarshal is a helper for unmarshaling as strictly as possible.
func strictUnmarshal(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}

	return validate(v)
}

// validate recursively validates a value.
func validate(v any) error {
	type validatable interface {
		Validate() error
	}

	if v, ok := any(v).(validatable); ok {
		return v.Validate()
	}
	return nil
}
