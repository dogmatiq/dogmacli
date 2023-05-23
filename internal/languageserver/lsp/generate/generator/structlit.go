package generator

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"

func (g *generator) VisitStructLit(n *model.StructLit) {
	name := nameOf(n)

	documentation(
		g,
		n.Documentation,
		"%s is a literal structure.",
		name,
	)

	g.emitStruct(
		name,
		nil, // embedded types
		n.Properties,
	)
}
