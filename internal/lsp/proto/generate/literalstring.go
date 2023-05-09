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
			Struct(),

		jen.Func().
			Params(jen.Id(name)).
			Id("MarshalJSON").
			Params().
			Params(
				jen.Index().Byte(),
				jen.Error(),
			).
			Block(
				jen.Return(
					jen.Id("marshal").Call(
						jen.Lit(t.LiteralString()),
					),
				),
			),

		jen.Func().
			Params(jen.Op("*").Id(name)).
			Id("UnmarshalJSON").
			Params(
				jen.Id("data").
					Index().
					Byte(),
			).
			Params(
				jen.Error(),
			).
			Block(
				jen.Var().Id("value").String(),
				jen.If(
					jen.Err().Op(":=").Id("unmarshal").Call(
						jen.Id("data"),
						jen.Op("&").Id("value"),
					),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Return(
						jen.Err(),
					),
				),

				jen.Line().
					If(
						jen.Id("value").
							Op("!=").
							Lit(t.LiteralString()),
					).
					Block(
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

				jen.Line().
					Return(
						jen.Nil(),
					),
			),
	)
}
