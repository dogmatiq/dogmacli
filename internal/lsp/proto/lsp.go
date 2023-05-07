package proto

import (
	"bytes"
	"encoding/json"
	"net/url"
)

// URI is a string that represents a URI.
type URI struct {
	url.URL
}

// MarshalJSON marshals the URI to JSON.
func (u URI) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON unmarshals the URI from JSON.
func (u *URI) UnmarshalJSON(data []byte) error {
	var s string
	if err := unmarshal(data, &s); err != nil {
		return err
	}

	p, err := url.Parse(s)
	if err != nil {
		return err
	}

	u.URL = *p
	return nil
}

// DocumentURI is a string that represents a URI for a document.
type DocumentURI = URI

func unmarshal(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
