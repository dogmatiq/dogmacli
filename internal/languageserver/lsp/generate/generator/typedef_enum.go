package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/jenx"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g typeDefGen) Enum(d model.Enum) {
	documentation(g, d.Documentation)
	g.
		Type().
		Id(exported(d.TypeName)).
		Add(typeExpr(d.Type))

	g.Line()
	g.enumConstants(d)

	g.Line()
	g.enumValidateMethod(d)

	g.Line()
	g.enumStringMethod(d)
}

func (g typeDefGen) enumConstants(d model.Enum) {
	g.
		Const().
		DefsFunc(func(g *jen.Group) {
			for _, m := range d.Members {
				documentation(g, m.Documentation)
				g.
					Id(exported(d.TypeName, m.Name)).
					Id(exported(d.TypeName)).
					Op("=").
					Lit(m.Value)
			}
		})
}

func (g typeDefGen) enumValidateMethod(d model.Enum) {
	g.Comment("Validate returns an error if x is invalid.")
	g.
		Func().
		Params(
			jen.Id("x").Id(exported(d.TypeName)),
		).
		Id("Validate").
		Params().
		Params(
			jen.Error(),
		).
		BlockFunc(func(g *jen.Group) {
			if !d.SupportsCustomValues {
				g.
					Switch(jen.Id("x")).
					BlockFunc(func(g *jen.Group) {
						for _, m := range d.Members {
							g.Case(
								jen.Id(exported(d.TypeName, m.Name)),
							)
						}

						g.
							Default().
							Block(
								jen.Return(
									jenx.
										Errorf(
											fmt.Sprintf(
												"invalid %s: %%v",
												exported(d.TypeName),
											),
											jen.Id("x"),
										),
								),
							)

					})
			}

			g.Return(jen.Nil())
		})
}

func (g typeDefGen) enumStringMethod(d model.Enum) {
	g.Comment("String returns the string representation of x.")
	g.
		Func().
		Params(
			jen.Id("x").Id(exported(d.TypeName)),
		).
		Id("String").
		Params().
		Params(
			jen.String(),
		).
		BlockFunc(func(g *jen.Group) {
			g.
				Switch(jen.Id("x")).
				BlockFunc(func(g *jen.Group) {
					for _, m := range d.Members {
						g.
							Case(
								jen.Id(exported(d.TypeName, m.Name)),
							).
							Block(
								jen.Return(
									jenx.Litf(
										"%s(%s)",
										exported(d.TypeName),
										exported(m.Name),
									),
								),
							)
					}

					tag := "invalid"
					if d.SupportsCustomValues {
						tag = "custom"
					}

					g.
						Default().
						Block(
							jen.Return(
								jenx.Sprintf(
									fmt.Sprintf(
										"%s(%%v, %s)",
										exported(d.TypeName),
										tag,
									),
									typeExpr(d.Type).Call(
										jen.Id("x"),
									),
								),
							),
						)

				})
		})
}
