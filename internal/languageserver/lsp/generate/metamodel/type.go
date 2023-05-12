package metamodel

import (
	"encoding/json"
)

// Type is an interface for a type.
type Type interface {
	acceptVisitor(TypeVisitor)
}

// VisitType dispatches to the method on v that corresponds to t's type.
func VisitType(t Type, v TypeVisitor) {
	t.acceptVisitor(v)
}

// TypeVisitor provides type-specific logic.
type TypeVisitor interface {
	VisitBool(BoolType)
	VisitDecimal(DecimalType)
	VisitString(StringType)
	VisitInteger(IntegerType)
	VisitUInteger(UIntegerType)
	VisitDocumentURI(DocumentURIType)
	VisitURI(URIType)
	VisitNull(NullType)
	VisitReference(ReferenceType)
	VisitArray(ArrayType)
	VisitMap(MapType)
	VisitLiteral(LiteralType)
	VisitStringLiteral(StringLiteralType)
	VisitAnd(AndType)
	VisitOr(OrType)
	VisitTuple(TupleType)
}

// BoolType is the "bool" base type.
type BoolType struct{}

// DecimalType is the "decimal" base type.
type DecimalType struct{}

// StringType is the "string" base type.
type StringType struct{}

// IntegerType is the "integer" base type.
type IntegerType struct{}

// UIntegerType is the "uinteger" base type.
type UIntegerType struct{}

// DocumentURIType is the "DocumentUri" base type.
type DocumentURIType struct{}

// URIType is the "URI" base type.
type URIType struct{}

// NullType is the "null" base type.
type NullType struct{}

// ReferenceType is a reference to a named type definition.
type ReferenceType struct {
	Target NamedType
}

// ArrayType is an array type.
type ArrayType struct {
	Element Type
}

// MapType is a map type.
type MapType struct {
	Key   Type
	Value Type
}

// LiteralType is a literal (anonymous) struct.
type LiteralType struct {
	Properties []StructureProperty
}

// StringLiteralType is a string type that must have a specific value. They are
// used as a discriminator in some union types.
type StringLiteralType struct {
	Value string
}

// AndType is a union type.
type AndType struct {
	Types []Type
}

// OrType is a union type.
type OrType struct {
	Types []Type
}

// TupleType is a tuple of values of (potentially) different types.
type TupleType struct {
	Types []Type
}

func newTypes(
	named map[string]NamedType,
	j []*typeJSON,
) []Type {
	var types []Type
	for _, t := range j {
		types = append(types, newType(named, t))
	}
	return types
}

func newType(
	named map[string]NamedType,
	j *typeJSON,
) Type {
	if j == nil {
		return nil
	}

	switch j.Kind {
	case "base":
		return newBaseType(j)
	case "reference":
		return ReferenceType{named[j.Name]}
	case "array":
		return ArrayType{newType(named, j.Element)}
	case "map":
		return newMapType(named, j)
	case "literal":
		return newLiteralType(named, j)
	case "stringLiteral":
		return newStringLiteralType(j)
	case "and":
		return AndType{newTypes(named, j.Items)}
	case "or":
		return OrType{newTypes(named, j.Items)}
	case "tuple":
		return TupleType{newTypes(named, j.Items)}
	default:
		panic("unsupported kind: " + j.Kind)
	}
}

func newBaseType(j *typeJSON) Type {
	switch j.Name {
	case "boolean":
		return BoolType{}
	case "decimal":
		return DecimalType{}
	case "string":
		return StringType{}
	case "integer":
		return IntegerType{}
	case "uinteger":
		return UIntegerType{}
	case "DocumentUri":
		return DocumentURIType{}
	case "URI":
		return URIType{}
	case "null":
		return NullType{}
	default:
		panic("unsupported base type: " + j.Name)
	}
}

func newMapType(
	named map[string]NamedType,
	t *typeJSON,
) MapType {
	var v typeJSON
	if err := json.Unmarshal(t.Value, &v); err != nil {
		panic(err)
	}

	return MapType{
		Key:   newType(named, t.Key),
		Value: newType(named, &v),
	}
}

func newLiteralType(
	named map[string]NamedType,
	t *typeJSON,
) LiteralType {
	var v struct {
		Properties []structurePropertyJSON `json:"properties"`
	}
	if err := json.Unmarshal(t.Value, &v); err != nil {
		panic(err)
	}

	return LiteralType{
		Properties: newStructureProperties(named, v.Properties),
	}
}

func newStringLiteralType(j *typeJSON) StringLiteralType {
	var v string
	if err := json.Unmarshal(j.Value, &v); err != nil {
		panic(err)
	}

	return StringLiteralType{
		Value: v,
	}
}

func (t BoolType) acceptVisitor(v TypeVisitor)          { v.VisitBool(t) }
func (t DecimalType) acceptVisitor(v TypeVisitor)       { v.VisitDecimal(t) }
func (t StringType) acceptVisitor(v TypeVisitor)        { v.VisitString(t) }
func (t IntegerType) acceptVisitor(v TypeVisitor)       { v.VisitInteger(t) }
func (t UIntegerType) acceptVisitor(v TypeVisitor)      { v.VisitUInteger(t) }
func (t DocumentURIType) acceptVisitor(v TypeVisitor)   { v.VisitDocumentURI(t) }
func (t URIType) acceptVisitor(v TypeVisitor)           { v.VisitURI(t) }
func (t NullType) acceptVisitor(v TypeVisitor)          { v.VisitNull(t) }
func (t ReferenceType) acceptVisitor(v TypeVisitor)     { v.VisitReference(t) }
func (t ArrayType) acceptVisitor(v TypeVisitor)         { v.VisitArray(t) }
func (t MapType) acceptVisitor(v TypeVisitor)           { v.VisitMap(t) }
func (t LiteralType) acceptVisitor(v TypeVisitor)       { v.VisitLiteral(t) }
func (t StringLiteralType) acceptVisitor(v TypeVisitor) { v.VisitStringLiteral(t) }
func (t AndType) acceptVisitor(v TypeVisitor)           { v.VisitAnd(t) }
func (t OrType) acceptVisitor(v TypeVisitor)            { v.VisitOr(t) }
func (t TupleType) acceptVisitor(v TypeVisitor)         { v.VisitTuple(t) }

type typeJSON struct {
	Kind    string          `json:"kind"`
	Name    string          `json:"name"`
	Items   []*typeJSON     `json:"items"`
	Element *typeJSON       `json:"element"`
	Key     *typeJSON       `json:"key"`
	Value   json.RawMessage `json:"value"`
}
