package lowlevel

import (
	"bytes"
	_ "embed"
	"encoding/json"
)

//go:embed metamodel-3.17.0.json
var data []byte

// Root returns the root node of the meta-model.
func Root() Model {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()

	var m Model
	if err := d.Decode(&m); err != nil {
		panic(err)
	}

	return m
}

// Documentation is a container for documentation-related meta-data.
type Documentation struct {
	Text               string `json:"documentation"`
	SinceVersion       string `json:"since"`
	DeprecationMessage string `json:"deprecated"`
}

// Model is the root of the meta-model.
type Model struct {
	MetaData struct {
		Version string `json:"version"`
	} `json:"metaData"`

	Requests      []Request      `json:"requests"`
	Notifications []Notification `json:"notifications"`
	Structs       []Struct       `json:"structures"`
	Enums         []Enum         `json:"enumerations"`
	Aliases       []Alias        `json:"typeAliases"`
}

// Request describes a JSON-RPC method that has a response.
type Request struct {
	Documentation

	Method              string `json:"method"`
	Direction           string `json:"messageDirection"`
	Params              Type   `json:"params"`
	Result              Type   `json:"result"`
	PartialResult       Type   `json:"partialResult"`
	ErrorData           Type   `json:"errorData"`
	RegistrationMethod  string `json:"registrationMethod"`
	RegistrationOptions Type   `json:"registrationOptions"`
}

// Notification describes a JSON-RPC method that does not have a response.
type Notification struct {
	Documentation

	Method              string `json:"method"`
	Direction           string `json:"messageDirection"`
	Params              Type   `json:"params"`
	RegistrationMethod  string `json:"registrationMethod"`
	RegistrationOptions Type   `json:"registrationOptions"`
}

// Struct describes a named structure type.
type Struct struct {
	Documentation

	Name       string     `json:"name"`
	Extends    []Type     `json:"extends"`
	Mixins     []Type     `json:"mixins"`
	Properties []Property `json:"properties"`
}

// Property describes a property within a structure.
type Property struct {
	Documentation

	Name     string `json:"name"`
	Optional bool   `json:"optional"`
	Type     Type   `json:"type"`
}

// Alias describes a named type alias.
type Alias struct {
	Documentation

	Name string `json:"name"`
	Type Type   `json:"type"`
}

// Enum describes a named enumeration type.
type Enum struct {
	Documentation

	Name                 string       `json:"name"`
	Type                 Type         `json:"type"`
	Members              []EnumMember `json:"values"`
	SupportsCustomValues bool         `json:"supportsCustomValues"`
}

// EnumMember describes a value within an enumeration.
type EnumMember struct {
	Documentation

	Name  string `json:"name"`
	Value any    `json:"value"`
}

// Type describes the type of a value.
type Type struct {
	Kind             Kind             `json:"kind"`
	Name             string           `json:"name"`
	Types            []Type           `json:"items"`
	ArrayElement     *Type            `json:"element"`
	MapKey           *Type            `json:"key"`
	MapValue         *Type            `json:"-"`
	StringLiteral    string           `json:"-"`
	StructureLiteral StructureLiteral `json:"-"`
	RawValue         json.RawMessage  `json:"value"`
}

// UnmarshalJSON unmarshals the JSON representation of the type.
func (t *Type) UnmarshalJSON(data []byte) error {
	type preventRecursion Type
	if err := json.Unmarshal(data, (*preventRecursion)(t)); err != nil {
		return err
	}

	switch t.Kind {
	case Map:
		return json.Unmarshal(t.RawValue, &t.MapValue)
	case Literal:
		return json.Unmarshal(t.RawValue, &t.StructureLiteral)
	case StringLiteral:
		return json.Unmarshal(t.RawValue, &t.StringLiteral)
	}

	return nil
}

// StructureLiteral describes an anonymous structure.
type StructureLiteral struct {
	Documentation

	Properties []Property `json:"properties"`
}

// Kind is an enumeration of the kinds of types.
type Kind string

const (
	// Base indicates that the Type refers to one of the base types.
	Base Kind = "base"

	// Reference indicates that the Type refers to a named type definition.
	Reference Kind = "reference"

	// Array indicates that the Type is an array of some other type.
	Array Kind = "array"

	// Map indicates that the Type is a map of one type to another.
	Map Kind = "map"

	// Literal indicates that the Type is a literal (anonymous) structure.
	Literal Kind = "literal"

	// StringLiteral indicates that the Type is a string with a specific value.
	StringLiteral Kind = "stringLiteral"

	// And indicates that the Type is an intersection of other types.
	And Kind = "and"

	// Or indicates that the Type is a union of other types.
	Or Kind = "or"

	// Tuple indicates that the Type is an n-tuple of other types.
	Tuple Kind = "tuple"
)

// BaseType is an enumeration of the base types.
type BaseType string

const (
	// Boolean is the base type for a boolean.
	Boolean BaseType = "boolean"

	// Decimal is the base type for a floating point number.
	Decimal BaseType = "decimal"

	// String is the base type for a string.
	String BaseType = "string"

	// Integer is the base type for a 32-bit signed integer.
	Integer BaseType = "integer"

	// UInteger is the base type for a 32-bit unsigned integer.
	UInteger BaseType = "uinteger"

	// DocumentURI is the base type for a URI that refers to a document.
	DocumentURI BaseType = "DocumentUri"

	// URI is the base type for a URI that does not refer to a document.
	URI BaseType = "URI"

	// Null is the base type for NULL.
	Null BaseType = "null"
)
