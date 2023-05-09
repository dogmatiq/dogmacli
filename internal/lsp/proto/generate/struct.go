package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateStructs(gen *jen.File) {
	generateBanner(gen, "STRUCTURES")

	for _, m := range g.root.Structures {
		g.generateStruct(gen, m)
		g.flushPending(gen)
	}
}

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

	name := normalizeName(m.Name)

	gen.Type().
		Id(name).
		StructFunc(func(gen *jen.Group) {
			for _, p := range m.Embeds() {
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

	expr := g.typeExpr(m.Type)
	tag := m.Name

	if m.Optional {
		tag += ",omitempty"

		if !g.isOmittable(m.Type) {
			expr = jen.
				Op("*").
				Add(expr)
		}
	}

	gen.Id(normalizeName(m.Name)).
		Add(expr).
		Tag(map[string]string{
			"json": tag,
		})
}
