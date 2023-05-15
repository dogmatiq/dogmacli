package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g typeDefGen) Alias(d model.Alias) {
	documentation(g, d.Documentation)
	g.
		Type().
		Id(exported(d.TypeName)).
		Op("=").
		Add(typeExpr(d.Type))
}
