package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeLit) StructLit(t model.StructLit) {
	g.File.
		Type().
		Id(g.Name).
		StructFunc(func(grp *jen.Group) {
			for _, p := range t.Properties {
				g.structProperty(grp, p)
			}
		})
}
