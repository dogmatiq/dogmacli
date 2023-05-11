package main

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) literalStringTypeExpr(t *metamodel.Type) jen.Code {
	name := normalizeUnexportedName(t.LiteralString() + "Literal")
	g.generateLiteralString(name, t)
	return jen.Id(name)
}

func (g *generator) generateLiteralString(name string, t *metamodel.Type) {
	g.pending = append(
		g.pending,
		jen.Commentf("%s is a type that must be represented as the JSON-string %q.", name, t.LiteralString()),
		jen.Type().
			Id(name).
			String(),

		jen.Func().
			Params(jen.Id("x").Id(name)).
			Id("Validate").
			Params().
			Params(
				jen.Error(),
			).
			Block(
				jen.If(
					jen.Id("x").
						Op("!=").
						Lit(t.LiteralString()),
				).Block(
					jen.Return(
						jen.Qual("errors", "New").
							Call(
								jen.Lit(
									fmt.Sprintf(
										"value must be %q",
										t.LiteralString(),
									),
								),
							),
					),
				),
				jen.Return(jen.Nil()),
			),
	)
}
