package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/generate/metamodel"
)

func (g *generator) literalStringRef(t *metamodel.Type) jen.Code {
	name := t.LiteralString() + "Literal"

	if _, ok := g.names[name]; ok {
		return jen.Id(name)
	}

	dataVar := t.LiteralString() + "JSON"

	g.pending = append(
		g.pending,
		jen.Commentf("%s is a type that must be represented as the JSON-string %q.", name, t.LiteralString()),
		jen.Type().
			Id(name).
			Struct(),

		jen.Var().
			Id(dataVar).
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
					jen.Id(dataVar),
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
							jen.Id(dataVar),
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
							jen.Id(dataVar),
						),
				),
			),
	)

	return jen.Id(name)
}
