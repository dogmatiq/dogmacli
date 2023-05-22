package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *generator) VisitEnum(n *model.Enum) {
	name := nameOf(n)
	underlying := nameOf(n.UnderlyingType)

	documentation(
		g,
		n.Documentation(),
		"%s is an enumeration of %s values.",
		name,
		underlying,
	)

	g.
		Type().
		Id(name).
		Id(underlying)

	g.
		Const().
		DefsFunc(
			g.withGroup(func() {
				for _, m := range n.Members {
					model.VisitNode(m, g)
				}
			}),
		)
}

func (g *generator) VisitEnumMember(n *model.EnumMember) {
	enum := nameOf(n.Parent())
	name := nameOf(n)

	documentation(
		g,
		n.Documentation,
		"%s is a member of the %s enumeration.",
		name,
		enum,
	)

	g.
		Id(name).
		Id(enum).
		Op("=").
		Lit(n.Value).
		Line()
}
