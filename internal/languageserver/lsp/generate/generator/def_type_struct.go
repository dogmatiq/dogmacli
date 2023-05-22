package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *generator) VisitStruct(n *model.Struct) {
	name := nameOf(n)

	documentation(
		g,
		n.Documentation(),
		"%s is a named structure definition.",
		name,
	)

	g.
		Type().
		Id(name).
		StructFunc(
			g.withGroup(func() {
				for _, t := range n.EmbeddedTypes {
					g.Id(nameOf(t))
				}

				if len(n.EmbeddedTypes) > 0 && len(n.Properties) > 0 {
					g.Line()
				}

				for _, p := range n.Properties {
					model.VisitNode(p, g)
				}
			}),
		)
}

func (g *generator) VisitProperty(n *model.Property) {
	name := nameOf(n)
	typ, ok := tryNameOf(n.Type)
	if !ok {
		return
	}

	documentation(
		g,
		n.Documentation,
		"",
	)

	expr := jen.Id(typ)
	if n.Optional {
		expr = jen.Id("Optional").Types(expr)
	}

	g.
		Id(name).
		Add(expr)
}
