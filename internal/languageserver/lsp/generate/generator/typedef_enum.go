package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/jenx"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Enum(d model.Enum) {
	documentation(
		g.File,
		d.Documentation,
		"Generated from the LSP '%s' enumeration.",
		d.TypeName,
	)

	g.emitEnumType(d)

	g.File.Line()
	g.emitEnumConstants(d)

	if !d.SupportsCustomValues {
		g.File.Line()
		g.emitEnumUnmarshalMethod(d)
	}
}

func (g *Generator) emitEnumType(d model.Enum) {
	info := g.typeInfoForDef(d)
	underlying := g.typeInfo(d.Type)

	g.File.
		Type().
		Add(info.TypeExpr()).
		Add(underlying.TypeExpr())
}

func (g *Generator) emitEnumConstants(d model.Enum) {
	info := g.typeInfoForDef(d)

	g.File.
		Const().
		DefsFunc(func(grp *jen.Group) {
			for _, m := range d.Members {
				documentation(grp, m.Documentation, "")
				grp.
					Id(identifier(*info.Name, m.Name)).
					Add(info.TypeExpr()).
					Op("=").
					Lit(m.Value)
			}
		})
}

func (g *Generator) emitEnumUnmarshalMethod(d model.Enum) {
	info := g.typeInfoForDef(d)

	g.File.
		Func().
		Params(
			jen.Id("x").Op("*").Add(info.TypeExpr()),
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
									jen.Op("*").Add(info.TypeExpr()),
								).
								Call(
									jen.Id("x"),
								),
						),
					jen.Err().Op("!=").Nil(),
				).
				Block(
					jen.Return(
						jenx.
							Errorf(
								fmt.Sprintf(
									"%s: %%w",
									*info.Name,
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
							jen.Id(identifier(*info.Name, m.Name)),
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
											*info.Name,
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
