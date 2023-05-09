package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) literalStringTypeExpr(t *metamodel.Type) jen.Code {
	name := normalizeUnexportedName(t.LiteralString() + "Literal")
	g.generateLiteralString(name, t)
	return jen.Id(name)
}

func (g *generator) generateLiteralString(name string, t *metamodel.Type) {
	data := t.LiteralString() + "JSON"

	g.pending = append(
		g.pending,
		jen.Commentf("%s is a type that must be represented as the JSON-string %q.", name, t.LiteralString()),
		jen.Type().
			Id(name).
			Struct(),

		jen.Var().
			Id(data).
			Op("=").
			Index().
			Byte().
			Call(jen.Lit(string(t.RawValue))),

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
					jen.Id(data),
					jen.Nil(),
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
				jen.If(
					jen.Qual("bytes", "Equal").
						Call(
							jen.Id("data"),
							jen.Id(data),
						),
				).Block(
					jen.Return(
						jen.Nil(),
					),
				),
				jen.Return(
					jen.Qual("fmt", "Errorf").
						Call(
							jen.Lit("unexpected JSON (%s), expected %s"),
							jen.Id("data"),
							jen.Id(data),
						),
				),
			),
	)
}
