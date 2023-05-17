package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Struct describes a named structure type.
type Struct struct {
	TypeName      string
	Documentation Documentation
	Embedded      []Type
	Properties    []Property
}

// Property describes a field within a structure.
type Property struct {
	Name          string
	Documentation Documentation
	Optional      bool
	Type          Type
}

// Name returns the type name.
func (d Struct) Name() string {
	return d.TypeName
}

func (b *builder) structure(in lowlevel.Struct, out *Struct) {
	out.TypeName = in.Name
	out.Documentation = in.Documentation
	out.Properties = b.properties(in.Properties)

	for _, t := range in.Extends {
		out.Embedded = append(
			out.Embedded,
			b.typeRef(t),
		)
	}

	for _, t := range in.Mixins {
		out.Embedded = append(
			out.Embedded,
			b.typeRef(t),
		)
	}
}

func (b *builder) properties(
	in []lowlevel.Property,
) (out []Property) {
	for _, p := range in {
		out = append(
			out,
			Property{
				Name:          p.Name,
				Documentation: p.Documentation,
				Optional:      p.Optional,
				Type:          b.typeRef(p.Type),
			},
		)
	}

	return out
}

func (d Struct) accept(v TypeDefVisitor) { v.Struct(d) }
func (v *typeDefX[T]) Struct(d Struct)   { v.V = v.X.Struct(d) }
