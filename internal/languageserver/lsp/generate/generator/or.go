package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *generator) VisitOr(n *model.Or) {
	name := nameOf(n)

	g.Commentf("%s is a union (aka 'or') of several other types.", name)
	g.Comment("")
	g.Comment("It may be one of the following types:")

	for _, t := range n.Types {
		if name, ok := tryNameOf(t); ok {
			g.Commentf("  - %s", name)
		}
	}

	g.
		Type().
		Id(name).
		Interface(
			jen.Id("is" + name).Params(),
		)
}
