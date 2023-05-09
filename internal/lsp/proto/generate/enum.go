package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateEnums(gen *jen.File) {
	generateBanner(gen, "ENUMERATIONS")

	for _, m := range g.root.Enumerations {
		g.generateEnum(gen, m)
		g.flushPending(gen)
	}
}

func (g *generator) generateEnum(
	gen *jen.File,
	m metamodel.Enumeration,
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
		Add(g.typeExpr(m.Type))

	gen.Const().
		DefsFunc(func(gen *jen.Group) {
			for _, member := range m.Members {
				generateDocs(gen, member.Documentation)

				value := member.Value
				if v, ok := value.(float64); ok {
					value = int(v)
				}

				gen.Id(normalizeName(m.Name) + normalizeName(member.Name)).
					Id(normalizeName(m.Name)).
					Op("=").
					Lit(value)
			}
		})
}
