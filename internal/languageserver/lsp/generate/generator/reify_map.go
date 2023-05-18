package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) Map(t model.Map) {
	documentation(
		g.File,
		model.Documentation{},
		"Generated from an LSP 'map' type.",
	)

	k := g.typeInfo(t.Key)
	v := g.typeInfo(t.Value)

	g.File.
		Type().
		Id(g.Name).
		Map(k.Expr()).
		Add(v.Expr())
}
