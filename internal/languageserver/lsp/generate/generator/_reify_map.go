package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) Map(t *model.Map) {
	documentation(
		g.File,
		model.Documentation{},
		"Generated from an LSP 'map' type.",
	)

	k := g.typeInfo(t.KeyType)
	v := g.typeInfo(t.ValueType)

	g.File.
		Type().
		Id(g.Name).
		Map(k.Expr()).
		Add(v.Expr())
}
