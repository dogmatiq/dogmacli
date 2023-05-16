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

	if !d.SupportsCustomValues {
		g.File.Line()
		g.enumUnmarshalMethod(d)
	}
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

func (g *Generator) enumUnmarshalMethod(d model.Enum) {
	g.File.
		Func().
		Params(
			jen.Id("x").Op("*").Id(identifier(d.TypeName)),
		).
		Id("UnmarshalJSON").
		Params(
			jen.Id("data").Index().Byte(),
		).
		Params(
			jen.Error(),
		).
		BlockFunc(func(grp *jen.Group) {
			grp.
				If(
					jen.
						Err().
						Op(":=").
						Qual("encoding/json", "Unmarshal").
						Call(
							jen.Id("data"),
							jen.
								Parens(
									jen.
										Op("*").
										Add(g.typeExpr(d.Type)),
								).
								Call(jen.Id("x")),
						),
					jen.Err().Op("!=").Nil(),
				).
				Block(
					jen.Return(
						jenx.
							Errorf(
								fmt.Sprintf(
									"%s: %%w",
									identifier(d.TypeName),
								),
								jen.Err(),
							),
					),
				)

			grp.
				Line().
				Switch(jen.Op("*").Id("x")).
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
											"%s: %%v is not a member of the enum",
											identifier(d.TypeName),
										),
										jen.Id("x"),
									),
							),
						)
				})

			grp.
				Line().
				Return(jen.Nil())
		})
}
