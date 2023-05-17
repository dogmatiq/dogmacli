package generator

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"

func (g *Generator) emitTypeDef(d model.TypeDef) {
	g.pushScope(d.Name())
	model.VisitTypeDef(d, &typeDef{g})
	g.popScope()
}

type typeDef struct{ *Generator }
