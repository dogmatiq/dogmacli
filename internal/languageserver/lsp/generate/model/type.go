package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Type is an interface for a type.
type Type interface {
	accept(TypeVisitor)
}

type (
	// Bool is the "bool" base type.
	Bool struct{}

	// Decimal is the "decimal" base type.
	Decimal struct{}

	// String is the "string" base type.
	String struct{}

	// Integer is the "integer" base type.
	Integer struct{}

	// UInteger is the "uinteger" base type.
	UInteger struct{}

	// DocumentURI is the "DocumentUri" base type.
	DocumentURI struct{}

	// URI is the "URI" base type.
	URI struct{}

	// Null is the "null" base type.
	Null struct{}

	// Reference is a reference to a named type.
	Reference struct {
		TargetName string
		Target     TypeDef
	}

	// StructLit is a literal (anonymous) struct.
	StructLit struct {
		Documentation Documentation
		Properties    []Property
	}

	// StringLit is a string that must have a specific value.
	StringLit struct{ Value string }

	// Array is an array of values of the same type.
	Array struct{ Element Type }

	// Map is a map of keys of one type to values of another type.
	Map struct{ Key, Value Type }

	// And is the intersection of multiple types.
	And struct{ Types []Type }

	// Or is the union of multiple types.
	Or struct{ Types []Type }

	// Tuple is an n-tuple of other types.
	Tuple struct{ Types []Type }
)

// TypeVisitor provides logic specific to each Type implementation.
type TypeVisitor interface {
	Bool()
	Decimal()
	String()
	Integer()
	UInteger()
	DocumentURI()
	URI()
	Null()
	Reference(Reference)
	Array(Array)
	Map(Map)
	And(And)
	Or(Or)
	Tuple(Tuple)
	StructLit(StructLit)
	StringLit(StringLit)
}

// VisitType dispatches to v based on the concrete type of t.
func VisitType(t Type, v TypeVisitor) {
	t.accept(v)
}

// TypeTransform produces a value of type T from a Type.
type TypeTransform[T any] interface {
	Bool() T
	Decimal() T
	String() T
	Integer() T
	UInteger() T
	DocumentURI() T
	URI() T
	Null() T
	Reference(Reference) T
	Array(Array) T
	Map(Map) T
	StructLit(StructLit) T
	StringLit(StringLit) T
	And(And) T
	Or(Or) T
	Tuple(Tuple) T
}

// TypeTo transforms t to a value of type T using x.
func TypeTo[T any](
	t Type,
	x TypeTransform[T],
) T {
	v := &typeX[T]{X: x}
	VisitType(t, v)
	return v.V
}

type typeX[T any] struct {
	X TypeTransform[T]
	V T
}

func (b *builder) typeRef(in lowlevel.Type) Type {
	var types []Type
	for _, t := range in.Types {
		types = append(types, b.typeRef(t))
	}

	switch in.Kind {
	case lowlevel.Base:
		return baseType(in)
	case lowlevel.Reference:
		return Reference{in.Name, b.typeDef(in.Name)}
	case lowlevel.Array:
		return Array{b.typeRef(*in.ArrayElement)}
	case lowlevel.Map:
		return Map{b.typeRef(*in.MapKey), b.typeRef(*in.MapValue)}
	case lowlevel.Literal:
		return StructLit{
			in.StructLit.Documentation,
			b.properties(in.StructLit.Properties),
		}
	case lowlevel.StringLit:
		return StringLit{in.StringLit}
	case lowlevel.And:
		return And{types}
	case lowlevel.Or:
		return Or{types}
	case lowlevel.Tuple:
		return Tuple{types}
	default:
		panic("unrecognized kind: " + in.Kind)
	}
}

func baseType(in lowlevel.Type) Type {
	switch lowlevel.BaseType(in.Name) {
	case lowlevel.Boolean:
		return Bool{}
	case lowlevel.Decimal:
		return Decimal{}
	case lowlevel.String:
		return String{}
	case lowlevel.Integer:
		return Integer{}
	case lowlevel.UInteger:
		return UInteger{}
	case lowlevel.DocumentURI:
		return DocumentURI{}
	case lowlevel.URI:
		return URI{}
	case lowlevel.Null:
		return Null{}
	default:
		panic("unrecognized base type: " + in.Name)
	}
}

func (t Bool) accept(v TypeVisitor)        { v.Bool() }
func (t Decimal) accept(v TypeVisitor)     { v.Decimal() }
func (t String) accept(v TypeVisitor)      { v.String() }
func (t Integer) accept(v TypeVisitor)     { v.Integer() }
func (t UInteger) accept(v TypeVisitor)    { v.UInteger() }
func (t DocumentURI) accept(v TypeVisitor) { v.DocumentURI() }
func (t URI) accept(v TypeVisitor)         { v.URI() }
func (t Null) accept(v TypeVisitor)        { v.Null() }
func (t Reference) accept(v TypeVisitor)   { v.Reference(t) }
func (t StructLit) accept(v TypeVisitor)   { v.StructLit(t) }
func (t StringLit) accept(v TypeVisitor)   { v.StringLit(t) }
func (t Array) accept(v TypeVisitor)       { v.Array(t) }
func (t Map) accept(v TypeVisitor)         { v.Map(t) }
func (t And) accept(v TypeVisitor)         { v.And(t) }
func (t Or) accept(v TypeVisitor)          { v.Or(t) }
func (t Tuple) accept(v TypeVisitor)       { v.Tuple(t) }

func (v *typeX[T]) Bool()                 { v.V = v.X.Bool() }
func (v *typeX[T]) Decimal()              { v.V = v.X.Decimal() }
func (v *typeX[T]) String()               { v.V = v.X.String() }
func (v *typeX[T]) Integer()              { v.V = v.X.Integer() }
func (v *typeX[T]) UInteger()             { v.V = v.X.UInteger() }
func (v *typeX[T]) DocumentURI()          { v.V = v.X.DocumentURI() }
func (v *typeX[T]) URI()                  { v.V = v.X.URI() }
func (v *typeX[T]) Null()                 { v.V = v.X.Null() }
func (v *typeX[T]) Reference(t Reference) { v.V = v.X.Reference(t) }
func (v *typeX[T]) StructLit(t StructLit) { v.V = v.X.StructLit(t) }
func (v *typeX[T]) StringLit(t StringLit) { v.V = v.X.StringLit(t) }
func (v *typeX[T]) Array(t Array)         { v.V = v.X.Array(t) }
func (v *typeX[T]) Map(t Map)             { v.V = v.X.Map(t) }
func (v *typeX[T]) And(t And)             { v.V = v.X.And(t) }
func (v *typeX[T]) Or(t Or)               { v.V = v.X.Or(t) }
func (v *typeX[T]) Tuple(t Tuple)         { v.V = v.X.Tuple(t) }
