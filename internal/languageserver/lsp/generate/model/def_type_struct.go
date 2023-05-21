package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Struct describes a named structure type.
type Struct struct {
	typeDefNode

	EmbeddedTypes []Type
	Properties    []*Property
}

// Property describes a field within a structure.
type Property struct {
	node

	Name          string
	Documentation Documentation
	Type          Type
	Optional      bool
}

func (b *builder) buildStruct(in lowlevel.Struct, out *Struct) {
	out.name = in.Name
	out.docs = in.Documentation

	for _, t := range in.Extends {
		out.EmbeddedTypes = append(out.EmbeddedTypes, b.buildType(t))
	}

	for _, t := range in.Mixins {
		out.EmbeddedTypes = append(out.EmbeddedTypes, b.buildType(t))
	}

	for _, p := range in.Properties {
		out.Properties = append(out.Properties, b.buildProperty(p))
	}
}

func (b *builder) buildProperty(in lowlevel.Property) *Property {
	return build(b, func(out *Property) {
		out.Name = in.Name
		out.Documentation = in.Documentation
		out.Type = b.buildType(in.Type)
		out.Optional = in.Optional
	})
}
