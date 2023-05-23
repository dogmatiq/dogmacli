package generator

import (
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
)

func (g *generator) VisitAnd(n *model.And) {
	name := nameOf(n)

	g.Commentf("%s is an intersection (aka 'and') of several other types.", name)

	g.emitStruct(
		name,
		n.Types,
		nil, // properties
	)
}
