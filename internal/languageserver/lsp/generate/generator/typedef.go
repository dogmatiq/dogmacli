package generator

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"

func (g *Generator) typeDef(d model.TypeDef) {
	g.pushScope(d.Name())
	defer g.popScope()

	model.VisitTypeDef(d, &typeDef{g})
}

type typeDef struct{ *Generator }
