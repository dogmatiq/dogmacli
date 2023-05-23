package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
)

func (g *generator) VisitTuple(n *model.Tuple) {
	name := nameOf(n)
	typ := nameOf(n.Types[0])

	for _, t := range n.Types[1:] {
		if nameOf(t) != typ {
			panic("tuple types must be homogeneous")
		}
	}

	g.
		Commentf("%s is a %d-tuple of %s.", name, len(n.Types), typ).
		Line().
		Type().
		Id(name).
		Index(jen.Lit(len(n.Types))).
		Id(typ)
}
