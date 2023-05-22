package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) And(t *model.And) {
	documentation(
		g.File,
		model.Documentation{},
		"Generated from an LSP 'and' type.",
	)
	g.emitStruct(
		g.Name,
		t.Types,
		nil, // properties
	)
}
