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

	g.emitStruct(
		name,
		n.EmbeddedTypes,
		n.Properties,
	)
}

func (g *generator) emitStruct(
	name string,
	embedded []model.Type,
	properties []*model.Property,
) {
	g.
		Type().
		Id(name).
		StructFunc(
			g.withGroup(func() {
				for _, t := range embedded {
					g.Id(nameOf(t))
				}

				if len(embedded) > 0 && len(properties) > 0 {
					g.Line()
				}

				for _, p := range properties {
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
