package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) Array(t *model.Array) {
	documentation(
		g.File,
		model.Documentation{},
		"Generated from an LSP 'array' type.",
	)

	e := g.typeInfo(t.Element)

	g.File.
		Type().
		Id(g.Name).
		Index().
		Add(e.Expr())

}
