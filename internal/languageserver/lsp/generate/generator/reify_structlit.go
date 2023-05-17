package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) StructLit(t model.StructLit) {
	g.emitStruct(
		g.Name,
		nil, // embedded
		t.Properties,
	)
}
