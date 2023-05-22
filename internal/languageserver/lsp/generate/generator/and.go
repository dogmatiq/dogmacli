package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *generator) VisitAnd(n *model.And) {
	name := nameOf(n)

	g.
		Commentf("%s is an intersection (aka 'and') of several other types.", name).
		Line().
		Type().
		Id(name).
		StructFunc(func(g *jen.Group) {
			for _, t := range n.Types {
				g.Id(nameOf(t))
			}
		})
}
