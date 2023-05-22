package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Alias(d *model.Alias) {
	info := g.typeInfoForDef(d)
	underlying := g.typeInfo(d.UnderlyingType)

	documentation(
		g.File,
		d.Documentation(),
		"Generated from the LSP '%s' type alias.",
		d.Name(),
	)

	if underlying.IsReified() {
		g.emitReifiedType(info.Name, d.UnderlyingType)
	} else {
		g.File.
			Type().
			Add(info.Expr()).
			Add(underlying.Expr())
	}
}
