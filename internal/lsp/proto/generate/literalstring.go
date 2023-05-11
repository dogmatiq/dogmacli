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

		jen.Line().
			Func().
			Params(jen.Id("x").Id(name)).
			Id("MarshalJSON").
			Params().
			Params(
				jen.Index().Byte(),
				jen.Error(),
			).
			Block(
				jen.Return(
					jen.Qual("encoding/json", "Marshal").Call(
						jen.Lit(t.LiteralString()),
					),
				),
			),

		jen.Line().
			Func().
			Params(jen.Id("x").Op("*").Id(name)).
			Id("UnmarshalJSON").
			Params(
				jen.Id("data").Index().Byte(),
			).
			Params(
				jen.Error(),
			).
			Block(
				jen.Var().Id("v").String(),
				jen.
					If(
						jen.Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(
							jen.Id("data"),
							jen.Op("&").Id("v"),
						),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(jen.Err()),
					),
				jen.
					If(
						jen.Id("v").Op("!=").Lit(t.LiteralString()),
					).
					Block(
						jen.Return(
							jen.Qual("errors", "New").Call(
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
