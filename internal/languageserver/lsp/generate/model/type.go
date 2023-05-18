package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Type is an interface for a type.
type Type interface {
	Node
	isType()
}

type typ struct{ node }

func (typ) isType() {}

type (
	// Bool is the "bool" base type.
	Bool struct{ typ }

	// Decimal is the "decimal" base type.
	Decimal struct{ typ }

	// String is the "string" base type.
	String struct{ typ }

	// Integer is the "integer" base type.
	Integer struct{ typ }

	// UInteger is the "uinteger" base type.
	UInteger struct{ typ }

	// DocumentURI is the "DocumentUri" base type.
	DocumentURI struct{ typ }

	// URI is the "URI" base type.
	URI struct{ typ }

	// Null is the "null" base type.
	Null struct{ typ }

	// Reference is a reference to a named type.
	Reference struct {
		typ
		Target TypeDef
	}

	// StructLit is a literal (anonymous) struct.
	StructLit struct {
		typ
		Documentation Documentation
		Properties    []*Property
	}

	// StringLit is a string that must have a specific value.
	StringLit struct {
		typ
		Value string
	}

	// Array is an array of values of the same type.
	Array struct {
		typ
		Element Type
	}

	// Map is a map of keys of one type to values of another type.
	Map struct {
		typ
		Key, Value Type
	}

	// And is the intersection of multiple types.
	And struct {
		typ
		Types []Type
	}

	// Or is the union of multiple types.
	Or struct {
		typ
		Types []Type
	}

	// Tuple is an n-tuple of other types.
	Tuple struct {
		typ
		Types []Type
	}
)

func (b *builder) buildType(in lowlevel.Type) Type {
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
		return build(b, func(n *Bool) {})
	case lowlevel.Decimal:
		return build(b, func(n *Decimal) {})
	case lowlevel.String:
		return build(b, func(n *String) {})
	case lowlevel.Integer:
		return build(b, func(n *Integer) {})
	case lowlevel.UInteger:
		return build(b, func(n *UInteger) {})
	case lowlevel.DocumentURI:
		return build(b, func(n *DocumentURI) {})
	case lowlevel.URI:
		return build(b, func(n *URI) {})
	case lowlevel.Null:
		return build(b, func(n *Null) {})
	default:
		panic("unrecognized base type: " + in.Name)
	}
}

func (b *builder) buildReferenceType(in lowlevel.Type) Type {
	return build(b, func(n *Reference) {
		if d, ok := b.aliases[in.Name]; ok {
			n.Target = d
		}
		if d, ok := b.enums[in.Name]; ok {
			n.Target = d
		}
		if d, ok := b.structs[in.Name]; ok {
			n.Target = d
		}
		if n.Target == nil {
			panic("unrecognized type: " + in.Name)
		}
	})
}

func (b *builder) buildStructLitType(in lowlevel.Type) Type {
	return build(b, func(n *StructLit) {
		n.Documentation = in.StructLit.Documentation
		n.Properties = b.buildProperties(in.StructLit.Properties)
	})
}

func (b *builder) buildStringLitType(in lowlevel.Type) Type {
	return build(b, func(n *StringLit) {
		n.Value = in.StringLit
	})
}

func (b *builder) buildArrayType(in lowlevel.Type) Type {
	return build(b, func(n *Array) {
		n.Element = b.buildType(*in.ArrayElement)
	})
}

func (b *builder) buildMapType(in lowlevel.Type) Type {
	return build(b, func(n *Map) {
		n.Key = b.buildType(*in.MapKey)
		n.Value = b.buildType(*in.MapValue)
	})
}

func (b *builder) buildAndType(in lowlevel.Type) Type {
	return build(b, func(n *And) {
		for _, t := range in.Types {
			n.Types = append(n.Types, b.buildType(t))
		}
	})
}

func (b *builder) buildOrType(in lowlevel.Type) Type {
	return build(b, func(n *Or) {
		for _, t := range in.Types {
			n.Types = append(n.Types, b.buildType(t))
		}
	})
}

func (b *builder) buildTupleType(in lowlevel.Type) Type {
	return build(b, func(n *Tuple) {
		for _, t := range in.Types {
			n.Types = append(n.Types, b.buildType(t))
		}
	})
}

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Bool) AcceptVisitor(v Visitor) { v.VisitBool(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Decimal) AcceptVisitor(v Visitor) { v.VisitDecimal(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *String) AcceptVisitor(v Visitor) { v.VisitString(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Integer) AcceptVisitor(v Visitor) { v.VisitInteger(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *UInteger) AcceptVisitor(v Visitor) { v.VisitUInteger(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *DocumentURI) AcceptVisitor(v Visitor) { v.VisitDocumentURI(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *URI) AcceptVisitor(v Visitor) { v.VisitURI(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Null) AcceptVisitor(v Visitor) { v.VisitNull(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Reference) AcceptVisitor(v Visitor) { v.VisitReference(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *StructLit) AcceptVisitor(v Visitor) { v.VisitStructLit(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *StringLit) AcceptVisitor(v Visitor) { v.VisitStringLit(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Array) AcceptVisitor(v Visitor) { v.VisitArray(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Map) AcceptVisitor(v Visitor) { v.VisitMap(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *And) AcceptVisitor(v Visitor) { v.VisitAnd(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Or) AcceptVisitor(v Visitor) { v.VisitOr(t) }

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (t *Tuple) AcceptVisitor(v Visitor) { v.VisitTuple(t) }
