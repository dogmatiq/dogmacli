package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/jenx"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Enum(d model.Enum) {
	documentation(g.File, d.Documentation)
	g.File.
		Type().
		Id(identifier(d.TypeName)).
		Add(g.typeExpr(d.Type))

	g.File.Line()
	g.enumConstants(d)

	// g.File.Line()
	// g.enumValidateMethod(d)

	// g.File.Line()
	// g.enumStringMethod(d)
}

func (g *Generator) enumConstants(d model.Enum) {
	g.File.
		Const().
		DefsFunc(func(grp *jen.Group) {
			for _, m := range d.Members {
				documentation(grp, m.Documentation)
				grp.
					Id(identifier(d.TypeName, m.Name)).
					Id(identifier(d.TypeName)).
					Op("=").
					Lit(m.Value)
			}
		})
}

func (g *Generator) enumValidateMethod(d model.Enum) {
	g.File.
		Comment("Validate returns an error if x is invalid.").
		Func().
		Params(
			jen.Id("x").Id(identifier(d.TypeName)),
		).
		Id("Validate").
		Params().
		Params(
			jen.Error(),
		).
		BlockFunc(func(grp *jen.Group) {
			if !d.SupportsCustomValues {
				grp.
					Switch(jen.Id("x")).
					BlockFunc(func(grp *jen.Group) {
						for _, m := range d.Members {
							grp.Case(
								jen.Id(identifier(d.TypeName, m.Name)),
							)
						}

						grp.
							Default().
							Block(
								jen.Return(
									jenx.
										Errorf(
											fmt.Sprintf(
												"invalid %s: %%v",
												identifier(d.TypeName),
											),
											jen.Id("x"),
										),
								),
							)

					})
			}

			grp.Return(jen.Nil())
		})
}

func (g *Generator) enumStringMethod(d model.Enum) {
	g.File.
		Comment("String returns the string representation of x.").
		Func().
		Params(
			jen.Id("x").Id(identifier(d.TypeName)),
		).
		Id("String").
		Params().
		Params(
			jen.String(),
		).
		BlockFunc(func(grp *jen.Group) {
			grp.
				Switch(jen.Id("x")).
				BlockFunc(func(grp *jen.Group) {
					for _, m := range d.Members {
						grp.
							Case(
								jen.Id(identifier(d.TypeName, m.Name)),
							).
							Block(
								jen.Return(
									jenx.Litf(
										"%s(%s)",
										identifier(d.TypeName),
										identifier(m.Name),
									),
								),
							)
					}

					tag := "invalid"
					if d.SupportsCustomValues {
						tag = "custom"
					}

					grp.
						Default().
						Block(
							jen.Return(
								jenx.Sprintf(
									fmt.Sprintf(
										"%s(%%v, %s)",
										identifier(d.TypeName),
										tag,
									),
									g.typeExpr(d.Type).Call(
										jen.Id("x"),
									),
								),
							),
						)

				})
		})
}
