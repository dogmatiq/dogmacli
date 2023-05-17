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
			name := rune('A')

			for _, m := range t.Types {
				info := g.typeInfo(m)

				if info.HasGoType {
					grp.
						Id(string(name)).
						Id("Optional").Types(
						g.typeExpr(m),
					)
				}

				name++
			}
		})
}
