package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/metamodel"
)

// VisitEnumeration declares a Go enumeration type.
func (g declarationGenerator) VisitEnumeration(t metamodel.Enumeration) {
	documentation(
		g,
		t.Documentation,
	)

	typeName := exported(t.Name)

	g.
		Type().
		Id(typeName).
		Add(typeExpr(t.Type))

	g.
		Const().
		DefsFunc(
			func(g *jen.Group) {
				for _, m := range t.Members {
					name := exported(t.Name, m.Name)

					documentation(
						g,
						m.Documentation,
					)

					g.
						Id(name).
						Id(typeName).
						Op("=").
						Lit(m.Value)
				}
			},
		)
}
