package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Alias(d model.Alias) {
	i := g.typeInfo(d.Type)

	documentation(g.File, d.Documentation)

	if i.IsLiteral {
		g.typeLit(d.TypeName, d.Type)
	} else {
		g.File.
			Type().
			Id(identifier(d.TypeName)).
			Op("=").
			Add(g.typeExpr(d.Type))
	}
}
