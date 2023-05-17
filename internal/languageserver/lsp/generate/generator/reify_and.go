package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) And(t model.And) {
	g.File.
		Type().
		Id(g.Name).
		StructFunc(func(grp *jen.Group) {
		})
}
