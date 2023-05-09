package metamodel

import (
	_ "embed"
	"encoding/json"

	"golang.org/x/exp/slices"
)

// Root is the root of the model.
type Root struct {
	Requests      []Request      `json:"requests"`
	Notifications []Notification `json:"notifications"`
	Structures    []Structure    `json:"structures"`
	Enumerations  []Enumeration  `json:"enumerations"`
	TypeAliases   []TypeAlias    `json:"typeAliases"`
}

// Request defines a JSON-RPC call (request/response).
type Request struct {
	Documentation       string `json:"documentation"`
	Method              string `json:"method"`
	Direction           string `json:"messageDirection"`
	Params              *Type  `json:"params"`
	Result              *Type  `json:"result"`
	PartialResult       *Type  `json:"partialResult"`
	RegistrationMethod  string `json:"registrationMethod"`
	RegistrationOptions *Type  `json:"registrationOptions"`
}

// Notification defines a JSON-RPC notification.
type Notification struct {
	Documentation       string `json:"documentation"`
	Method              string `json:"method"`
	Direction           string `json:"messageDirection"`
	Params              *Type  `json:"params"`
	RegistrationMethod  string `json:"registrationMethod"`
	RegistrationOptions *Type  `json:"registrationOptions"`
}

// Type is a reference to a named type, or an inline anonymous type.
type Type struct {
	Kind         string          `json:"kind"`
	Name         string          `json:"name"`
	Items        []*Type         `json:"items"`
	ArrayElement *Type           `json:"element"`
	MapKey       *Type           `json:"key"`
	RawValue     json.RawMessage `json:"value"`
}

// IsNull returns true if the type is the null type.
func (t *Type) IsNull() bool {
	return t.Kind == "base" && t.Name == "null"
}

// MapValue unserializes RawValue for use as a map value.
func (t *Type) MapValue() *Type {
	if t.Kind != "map" {
		panic("not a map type")
	}

	var value Type

	if err := json.Unmarshal(t.RawValue, &value); err != nil {
		panic(err)
	}

	return &value
}

// LiteralString unserializes RawValue for use as a literal string type.
func (t *Type) LiteralString() string {
	if t.Kind != "stringLiteral" {
		panic("not a literal string type")
	}

	var value string

	if err := json.Unmarshal(t.RawValue, &value); err != nil {
		panic(err)
	}

	return value
}

// LiteralStructProperties unserializes RawValue for use as a literal type's
// properties.
func (t *Type) LiteralStructProperties() []Property {
	if t.Kind != "literal" {
		panic("not a literal type")
	}

	var value struct {
		Properties []Property `json:"properties"`
	}

	if err := json.Unmarshal(t.RawValue, &value); err != nil {
		panic(err)
	}

	return value.Properties
}

// Structure defines a named structure data type.
type Structure struct {
	Documentation string     `json:"documentation"`
	Name          string     `json:"name"`
	Properties    []Property `json:"properties"`
	Extends       []*Type    `json:"extends"`
	Mixins        []*Type    `json:"mixins"`
}

// Embeds returns the types that this structure embeds.
func (s Structure) Embeds() []*Type {
	return append(
		slices.Clone(s.Extends),
		s.Mixins...,
	)
}

// Property is a member of a structure.
type Property struct {
	Documentation string `json:"documentation"`
	Name          string `json:"name"`
	Type          *Type  `json:"type"`
	Optional      bool   `json:"optional"`
}

// Enumeration defines a named enumeration data type.
type Enumeration struct {
	Documentation string              `json:"documentation"`
	Name          string              `json:"name"`
	Type          *Type               `json:"type"`
	Members       []EnumerationMember `json:"values"`
}

// EnumerationMember is a member of an enumeration.
type EnumerationMember struct {
	Documentation string `json:"documentation"`
	Name          string `json:"name"`
	Value         any    `json:"value"`
}

// TypeAlias defines a named type alias.
type TypeAlias struct {
	Documentation string `json:"documentation"`
	Name          string `json:"name"`
	Type          *Type  `json:"type"`
}

//go:embed metamodel-3.17.0.json
var data []byte

// Get returns the root node of the model.
func Get() Root {
	var root Root

	if err := json.Unmarshal(data, &root); err != nil {
		panic(err)
	}

	slices.SortFunc(
		root.Requests,
		func(a, b Request) bool {
			return a.Method < b.Method
		},
	)

	slices.SortFunc(
		root.Structures,
		func(a, b Structure) bool {
			return a.Name < b.Name
		},
	)

	slices.SortFunc(
		root.Enumerations,
		func(a, b Enumeration) bool {
			return a.Name < b.Name
		},
	)

	slices.SortFunc(
		root.TypeAliases,
		func(a, b TypeAlias) bool {
			return a.Name < b.Name
		},
	)

	return root
}
