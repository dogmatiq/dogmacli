package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
)

func (g *generator) VisitBool(n *model.Bool) {
	name := nameOf(n)

	g.
		Commentf("%s is the LSP boolean type.", name).
		Line().
		Type().
		Id(name).
		Bool()

	g.
		Const().
		Defs(
			jen.Id("True").
				Id(name).
				Op("=").
				True(),
			jen.Id("False").
				Id(name).
				Op("=").
				False(),
		)

	g.
		Line().
		Func().
		Params(jen.Id("x").Id(name)).
		Id("isZero").
		Params().
		Bool().
		Block(
			jen.Return(jen.Id("x").Op("==").Id("False")),
		)
}

func (g *generator) VisitDecimal(n *model.Decimal) {
	name := nameOf(n)

	g.
		Commentf("%s is the LSP decimal type.", name).
		Line().
		Type().
		Id(name).
		Float64()

	g.
		Line().
		Func().
		Params(jen.Id("x").Id(name)).
		Id("isZero").
		Params().
		Bool().
		Block(
			jen.Return(jen.Id("x").Op("==").Lit(0)),
		)
}

func (g *generator) VisitString(n *model.String) {
	name := nameOf(n)

	g.
		Commentf("%s is the LSP string type.", name).
		Line().
		Type().
		Id(name).
		String()

	g.
		Line().
		Func().
		Params(jen.Id("x").Id(name)).
		Id("isZero").
		Params().
		Bool().
		Block(
			jen.Return(jen.Id("x").Op("==").Lit("")),
		)
}

func (g *generator) VisitInteger(n *model.Integer) {
	name := nameOf(n)

	g.
		Commentf("%s is the LSP signed integer type.", name).
		Line().
		Type().
		Id(name).
		Int32()

	g.
		Line().
		Func().
		Params(jen.Id("x").Id(name)).
		Id("isZero").
		Params().
		Bool().
		Block(
			jen.Return(jen.Id("x").Op("==").Lit(0)),
		)
}

func (g *generator) VisitUInteger(n *model.UInteger) {
	name := nameOf(n)

	g.
		Commentf("%s is the LSP unsigned integer type.", name).
		Line().
		Type().
		Id(name).
		Uint32()

	g.
		Line().
		Func().
		Params(jen.Id("x").Id(name)).
		Id("isZero").
		Params().
		Bool().
		Block(
			jen.Return(jen.Id("x").Op("==").Lit(0)),
		)
}

func (g *generator) VisitDocumentURI(n *model.DocumentURI) {
	name := nameOf(n)

	g.
		Commentf("%s is the URI of a document.", name).
		Line().
		Type().
		Id(name).
		Qual("net/url", "URL")

	g.
		Line().
		Func().
		Params(jen.Id("x").Id(name)).
		Id("isZero").
		Params().
		Bool().
		Block(
			jen.Return(
				jen.
					Id("x").Op("==").
					Parens(jen.Id(name).Values()),
			),
		)
}

func (g *generator) VisitURI(n *model.URI) {
	name := nameOf(n)

	g.
		Commentf("%s is the URI of some non-document resource.", name).
		Line().
		Type().
		Id(name).
		Qual("net/url", "URL")

	g.
		Line().
		Func().
		Params(jen.Id("x").Id(name)).
		Id("isZero").
		Params().
		Bool().
		Block(
			jen.Return(
				jen.
					Id("x").Op("==").
					Parens(jen.Id(name).Values()),
			),
		)
}

func (g *generator) VisitNull(n *model.Null) {}
