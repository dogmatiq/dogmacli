package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) StructLit(t model.StructLit) {
	documentation(
		g.File,
		model.Documentation{},
		"Generated from an LSP 'literal' type.",
	)
	g.emitStruct(
		g.Name,
		nil, // embedded
		t.Properties,
	)
}
