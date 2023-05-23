package lsp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func marshalProperty(
	w *bytes.Buffer,
	n *int,
	k string,
	v any,
) error {
	if v == nil {
		return nil
	}

	*n++
	if *n == 1 {
		w.WriteByte(',')
	}

	enc := json.NewEncoder(w)

	if err := enc.Encode(k); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}

	w.WriteByte(':')

	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}

	return nil
}

func marshalOptionalProperty[T any](
	w *bytes.Buffer,
	n *int,
	k string,
	v Optional[T],
) error {
	if v, ok := v.Get(); ok {
		return marshalProperty(w, n, k, v)
	}
	return nil
}

func unmarshalProperty(
	p map[string]json.RawMessage,
	k string,
	v any,
) error {
	data, ok := p[k]
	if !ok {
		return fmt.Errorf("%s: mandatory property is absent", k)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}

	return nil
}

func unmarshalOptionalProperty[T any](
	p map[string]json.RawMessage,
	k string,
	v *Optional[T],
) error {
	data, ok := p[k]
	if !ok {
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}

	*v = Optional[T]{
		value: &value,
	}

	return nil
}

func unmarshalPropertyUsing[T any](
	p map[string]json.RawMessage,
	k string,
	v *T,
	fn func([]byte, *T) error,
) error {
	data, ok := p[k]
	if !ok {
		return fmt.Errorf("%s: mandatory property is absent", k)
	}

	if err := fn(data, v); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}

	return nil
}

func unmarshalOptionalPropertyUsing[T any](
	p map[string]json.RawMessage,
	k string,
	v *Optional[T],
	fn func([]byte, *T) error,
) error {
	data, ok := p[k]
	if !ok {
		return nil
	}

	var value T
	if err := fn(data, &value); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}

	*v = Optional[T]{
		value: &value,
	}

	return nil
}

func unmarshalLiteralProperty(
	p map[string]json.RawMessage,
	k string,
	expect string,
) error {
	data, ok := p[k]
	if !ok {
		return fmt.Errorf("%s: mandatory property is absent", k)
	}

	var actual string
	if err := json.Unmarshal(data, &actual); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}

	if actual != expect {
		return fmt.Errorf("%s: expected %q, got %q", k, expect, actual)
	}

	return nil
}
