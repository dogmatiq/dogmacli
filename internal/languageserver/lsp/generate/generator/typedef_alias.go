package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Alias(d model.Alias) {
	i := g.typeInfo(d.Type)

	documentation(g.File, d.Documentation)

	if i.IsLiteral {
		g.emitReifiedType(d.TypeName, d.Type)
	} else {
		g.emitAliasType(d)
	}
}

func (g *Generator) emitAliasType(d model.Alias) {
	g.File.
		Type().
		Id(identifier(d.TypeName)).
		Op("=").
		Add(g.typeExpr(d.Type))
}
