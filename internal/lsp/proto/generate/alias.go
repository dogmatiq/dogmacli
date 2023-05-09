package main

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateAliases(gen *jen.File) {
	generateBanner(gen, "TYPE ALIASES")

	for _, m := range g.root.TypeAliases {
		if !strings.HasPrefix(m.Name, "LSP") {
			g.generateAlias(gen, m)
			g.flushPending(gen)
		}
	}
}

func (g *generator) generateAlias(
	gen *jen.File,
	m metamodel.TypeAlias,
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
		Op("=").
		Add(g.typeExpr(m.Type))
}
