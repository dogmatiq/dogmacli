package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) Tuple(t model.Tuple) {
	g.File.
		Type().
		Id(g.Name).
		StructFunc(func(grp *jen.Group) {
		})
}
