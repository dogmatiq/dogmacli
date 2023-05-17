package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) Or(t model.Or) {
	g.File.
		Type().
		Id(g.Name).
		StructFunc(func(grp *jen.Group) {
			for i, m := range t.Types {
				info := g.typeInfo(m)
				expr := g.typeExpr(m)
				name := string(rune('A' + i))

				if !info.IsNativelyOptional {
					expr = jen.Id("Optional").Types(expr)
				}

				grp.
					Id(name).
					Add(expr)
			}
		})
}
