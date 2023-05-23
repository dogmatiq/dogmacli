package generator

import "github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"

func (g *generator) VisitArray(n *model.Array) {
	name := nameOf(n)
	elem := nameOf(n.ElementType)

	g.
		Commentf("%s is an array of %s elements.", name, elem).
		Line().
		Type().
		Id(name).
		Index().
		Id(elem)
}
