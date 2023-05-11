package proto

import (
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
	if err := json.Unmarshal(data, &s); err != nil {
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

// Array is an array of T.
type Array[T any] []T

// Validate validates each element of the array.
func (a Array[T]) Validate() error {
	for _, v := range a {
		if err := validate(v); err != nil {
			return err
		}
	}
	return nil
}

// Map is a map of K to V.
type Map[K comparable, V any] map[K]V

// Validate validates each key and value of the map.
func (m Map[K, V]) Validate() error {
	for k, v := range m {
		if err := validate(k); err != nil {
			return err
		}
		if err := validate(v); err != nil {
			return err
		}
	}
	return nil
}

// V returns a pointer to v.
func V[T any](v T) *T {
	return &v
}
