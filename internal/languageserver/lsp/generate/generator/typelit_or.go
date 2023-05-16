package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeLit) Or(t model.Or) {
	g.File.
		Type().
		Id(g.Name).
		StructFunc(func(grp *jen.Group) {
		})
}
