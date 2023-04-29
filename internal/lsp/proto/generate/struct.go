package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateStruct(
	gen *jen.File,
	m metamodel.Structure,
) {
	g.pushName(m.Name)
	defer g.popName()

	if m.Documentation == "" {
		gen.Line()
	} else {
		generateDocs(gen, m.Documentation)
	}

	gen.Type().
		Id(normalizeName(m.Name)).
		StructFunc(func(gen *jen.Group) {
			for _, p := range m.Extends {
				gen.Id(normalizeName(p.Name))
			}
			for _, p := range m.Mixins {
				gen.Id(normalizeName(p.Name))
			}
			for _, p := range m.Properties {
				g.generateStructProperty(gen, p)
			}
		})
}

func (g *generator) generateStructProperty(
	gen *jen.Group,
	m metamodel.Property,
) {
	g.pushName(m.Name)
	defer g.popName()

	generateDocs(gen, m.Documentation)

	ref := g.typeRef(m.Type)
	tag := m.Name

	if m.Optional {
		tag += ",omitempty"

		if !g.isOmittable(m.Type) {
			ref = jen.
				Op("*").
				Add(ref)
		}
	}

	gen.Id(normalizeName(m.Name)).
		Add(ref).
		Tag(map[string]string{
			"json": tag,
		})
}
