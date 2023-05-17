package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) And(t model.And) {
	g.emitStruct(
		g.Name,
		t.Types,
		nil, // properties
	)
}
