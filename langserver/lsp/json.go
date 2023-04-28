package lsp

import (
	"bytes"
	"encoding/json"
)

func unmarshal(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
