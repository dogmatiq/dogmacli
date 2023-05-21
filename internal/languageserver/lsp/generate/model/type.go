package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Type describes the type of some value.
type Type interface {
	Node

	IsAnonymous() bool
}

// typeNode provides implementation common to all types that implement Type.
type typeNode struct {
	node
	anon bool
}

func (n *typeNode) IsAnonymous() bool {
	return n.anon
}

// TypeVisitor is an interface for type specific logic.
type TypeVisitor interface {
	VisitBool(*Bool)
	VisitDecimal(*Decimal)
	VisitString(*String)
	VisitInteger(*Integer)
	VisitUInteger(*UInteger)
	VisitDocumentURI(*DocumentURI)
	VisitURI(*URI)
	VisitNull(*Null)
	VisitReference(*Reference)
	VisitArray(*Array)
	VisitMap(*Map)
	VisitAnd(*And)
	VisitOr(*Or)
	VisitTuple(*Tuple)
	VisitStructLit(*StructLit)
	VisitStringLit(*StringLit)
}

// VisitType dispatches to the appropriate method on the given visitor based on
// the concrete type of t.
func VisitType(t Type, v TypeVisitor) {
	type satisfy struct{ NodeVisitor }
	type visitor struct {
		// embed a NodeVisitor to satisfy the other methods of NodeVisitor. They
		// will panic if called.
		satisfy
		TypeVisitor
	}

	t.acceptVisitor(visitor{TypeVisitor: v})
}

type (
	// Bool is the "bool" base type.
	Bool struct{ typeNode }

	// Decimal is the "decimal" base type.
	Decimal struct{ typeNode }

	// String is the "string" base type.
	String struct{ typeNode }

	// Integer is the "integer" base type.
	Integer struct{ typeNode }

	// UInteger is the "uinteger" base type.
	UInteger struct{ typeNode }

	// DocumentURI is the "DocumentUri" base type.
	DocumentURI struct{ typeNode }

	// URI is the "URI" base type.
	URI struct{ typeNode }

	// Null is the "null" base type.
	Null struct{ typeNode }

	// Reference is a reference to a named type.
	Reference struct {
		typeNode
		Target TypeDef
	}

	// StructLit is a literal (anonymous) struct.
	StructLit struct {
		typeNode
		Documentation Documentation
		Properties    []*Property
	}

	// StringLit is a string that must have a specific value.
	StringLit struct {
		typeNode
		Value string
	}

	// Array is an array of values of the same type.
	Array struct {
		typeNode
		ElementType Type
	}

	// Map is a map of keys of one type to values of another type.
	Map struct {
		typeNode
		KeyType, ValueType Type
	}

	// And is the intersection of multiple types.
	And struct {
		typeNode
		Types []Type
	}

	// Or is the union of multiple types.
	Or struct {
		typeNode
		Types []Type
	}

	// Tuple is an n-tuple of other types.
	Tuple struct {
		typeNode
		Types []Type
	}
)

func (b *builder) buildType(in lowlevel.Type) (out Type) {
	defer func() {
		if out != nil {
			b.model.Types = append(b.model.Types, out)
		}
	}()

	switch in.Kind {
	case "":
		return nil
	case lowlevel.Base:
		return b.buildBaseType(in)
	case lowlevel.Reference:
		return b.buildReferenceType(in)
	case lowlevel.Literal:
		return b.buildStructLitType(in)
	case lowlevel.StringLiteral:
		return b.buildStringLitType(in)
	case lowlevel.Array:
		return b.buildArrayType(in)
	case lowlevel.Map:
		return b.buildMapType(in)
	case lowlevel.And:
		return b.buildAndType(in)
	case lowlevel.Or:
		return b.buildOrType(in)
	case lowlevel.Tuple:
		return b.buildTupleType(in)
	default:
		panic("unrecognized kind: " + in.Kind)
	}
}

func (b *builder) buildBaseType(in lowlevel.Type) Type {
	switch lowlevel.BaseType(in.Name) {
	case lowlevel.Boolean:
		return build(b, func(*Bool) {})
	case lowlevel.Decimal:
		return build(b, func(*Decimal) {})
	case lowlevel.String:
		return build(b, func(*String) {})
	case lowlevel.Integer:
		return build(b, func(*Integer) {})
	case lowlevel.UInteger:
		return build(b, func(*UInteger) {})
	case lowlevel.DocumentURI:
		return build(b, func(*DocumentURI) {})
	case lowlevel.URI:
		return build(b, func(*URI) {})
	case lowlevel.Null:
		return build(b, func(*Null) {})
	default:
		panic("unrecognized base type: " + in.Name)
	}
}

func (b *builder) buildReferenceType(in lowlevel.Type) Type {
	return build(b, func(out *Reference) {
		out.Target = b.model.Defs[in.Name].(TypeDef)
	})
}

func (b *builder) buildStructLitType(in lowlevel.Type) Type {
	return build(b, func(out *StructLit) {
		out.anon = true

		out.Documentation = in.StructLit.Documentation

		for _, p := range in.StructLit.Properties {
			out.Properties = append(out.Properties, b.buildProperty(p))
		}
	})
}

func (b *builder) buildStringLitType(in lowlevel.Type) Type {
	return build(b, func(out *StringLit) {
		out.anon = true

		out.Value = in.StringLit
	})
}

func (b *builder) buildArrayType(in lowlevel.Type) Type {
	return build(b, func(out *Array) {
		out.anon = true

		out.ElementType = b.buildType(*in.ArrayElement)
	})
}

func (b *builder) buildMapType(in lowlevel.Type) Type {
	return build(b, func(out *Map) {
		out.anon = true

		out.KeyType = b.buildType(*in.MapKey)
		out.ValueType = b.buildType(*in.MapValue)
	})
}

func (b *builder) buildAndType(in lowlevel.Type) Type {
	return build(b, func(out *And) {
		out.anon = true

		for _, t := range in.Types {
			out.Types = append(out.Types, b.buildType(t))
		}
	})
}

func (b *builder) buildOrType(in lowlevel.Type) Type {
	return build(b, func(out *Or) {
		out.anon = true

		for _, t := range in.Types {
			out.Types = append(out.Types, b.buildType(t))
		}
	})
}

func (b *builder) buildTupleType(in lowlevel.Type) Type {
	return build(b, func(out *Tuple) {
		out.anon = true

		for _, t := range in.Types {
			out.Types = append(out.Types, b.buildType(t))
		}
	})
}
