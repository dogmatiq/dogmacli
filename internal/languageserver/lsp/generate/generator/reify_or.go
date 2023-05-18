package generator

import (
	"reflect"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) Or(t model.Or) {
	documentation(
		g.File,
		model.Documentation{},
		"Generated from an LSP 'or' type.",
	)
	g.File.
		Type().
		Id(g.Name).
		StructFunc(func(grp *jen.Group) {
			for i, m := range t.Types {
				g.pushNestedScope(strconv.Itoa(i))

				info := g.typeInfo(m)
				if info.TypeKind != reflect.Invalid {
					grp.
						Id(info.NameHint).
						Add(info.TypeExpr())
				}

				g.popNestedScope()
			}
		})
}
