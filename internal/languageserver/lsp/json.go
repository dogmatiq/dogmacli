package lsp

import (
	"bytes"
	"encoding/json"
)

func marshalProperty(
	w *bytes.Buffer,
	n *int,
	k string,
	v any,
) error {
	*n++
	if *n == 1 {
		w.WriteByte(',')
	}

	enc := json.NewEncoder(w)

	if err := enc.Encode(k); err != nil {
		return err
	}

	w.WriteByte(':')

	return enc.Encode(v)
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
