package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Struct describes a named structure type.
type Struct struct {
	typeDef

	Embedded   []Type
	Properties []*Property
}

// Property describes a field within a structure.
type Property struct {
	node

	Name          string
	Documentation Documentation
	Optional      bool
	Type          Type
}

func (b *builder) buildStruct(in lowlevel.Struct, out *Struct) {
	out.Documentation = in.Documentation
	out.Properties = b.buildProperties(in.Properties)

	for _, t := range in.Extends {
		out.Embedded = append(
			out.Embedded,
			b.buildType(t),
		)
	}

	for _, t := range in.Mixins {
		out.Embedded = append(
			out.Embedded,
			b.buildType(t),
		)
	}
}

func (b *builder) buildProperties(in []lowlevel.Property) []*Property {
	var out []*Property

	for _, p := range in {
		out = append(
			out,
			build(b, func(n *Property) {
				n.Name = p.Name
				n.Documentation = p.Documentation
				n.Optional = p.Optional
				n.Type = b.buildType(p.Type)
			}),
		)
	}

	return out
}
