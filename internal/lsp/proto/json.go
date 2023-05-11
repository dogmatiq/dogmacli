package proto

import (
	"bytes"
	"encoding/json"
)

// structUnmarshal is a helper for unmarshaling as strictly as possible.
func strictUnmarshal[T any](data []byte, v *T) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}

	return validate(*v)
}

// validate recursively validates a value.
func validate(v any) error {
	type validatable interface {
		Validate() error
	}

	if v, ok := v.(validatable); ok {
		return v.Validate()
	}
	return nil
}
