package proto

import (
	"bytes"
	"encoding/json"
)

// URI is a string that represents a URI.
type URI = string

// DocumentURI is a string that represents a URI for a document.
type DocumentURI = string

func unmarshal(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
