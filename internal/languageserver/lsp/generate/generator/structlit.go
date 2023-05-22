package generator

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"

func (g *generator) VisitStructLit(n *model.StructLit) {
	name := nameOf(n)

	g.
		Commentf("%s is a literal structure.", name).
		Line().
		Type().
		Id(name).
		Struct()
}
