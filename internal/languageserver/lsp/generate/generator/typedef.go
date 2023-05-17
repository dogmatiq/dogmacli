package generator

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"

func (g *Generator) emitTypeDef(d model.TypeDef) {
	g.enterType(d.Name())
	model.VisitTypeDef(d, &typeDef{g})
	g.leaveType()
}

type typeDef struct{ *Generator }
