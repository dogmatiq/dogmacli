package generator

import (
	"reflect"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *reifyType) Or(t *model.Or) {
	tag := "is" + g.Name

	documentation(
		g.File,
		model.Documentation{},
		"Generated from an LSP 'or' type.",
	)
	g.File.
		Type().
		Id(g.Name).
		Interface(
			jen.
				Id(tag).
				Params(),
		)

	for i, m := range t.Types {
		g.pushNestedScope("Option" + strconv.Itoa(i))
		info := g.typeInfo(m)
		g.popNestedScope()

		if info.Kind == reflect.Invalid {
			continue
		}

		if info.Kind == reflect.Interface {
		} else {
			g.File.
				Func().
				Params(info.Expr()).
				Id(tag).
				Params().
				Block()
		}
	}
}
