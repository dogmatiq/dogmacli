package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/jenx"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Enum(d *model.Enum) {
	documentation(
		g.File,
		d.Documentation(),
		"Generated from the LSP '%s' enumeration.",
		d.Name(),
	)

	g.emitEnumType(d)

	g.File.Line()
	g.emitEnumConstants(d)

	if d.Strict {
		g.File.Line()
		g.emitEnumUnmarshalMethod(d)
	}
}

func enumMemberName(d *model.Enum, m *model.EnumMember) string {
	return normalize(m.Name, d.Name())
}

func (g *Generator) emitEnumType(d *model.Enum) {
	info := g.typeInfoForDef(d)
	underlying := g.typeInfo(d.UnderlyingType)

	g.File.
		Type().
		Add(info.Expr()).
		Add(underlying.Expr())
}

func (g *Generator) emitEnumConstants(d *model.Enum) {
	info := g.typeInfoForDef(d)

	g.File.
		Const().
		DefsFunc(func(grp *jen.Group) {
			for _, m := range d.Members {
				documentation(grp, m.Documentation, "")
				grp.
					Id(enumMemberName(d, m)).
					Add(info.Expr()).
					Op("=").
					Lit(m.Value)
			}
		})
}

func (g *Generator) emitEnumUnmarshalMethod(d *model.Enum) {
	info := g.typeInfoForDef(d)

	g.File.
		Func().
		Params(
			jen.Id("x").Op("*").Add(info.Expr()),
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
									jen.Op("*").Add(info.Expr()),
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
								fmt.Sprintf("%s: %%w", info.Name),
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
							jen.Id(enumMemberName(d, m)),
						)
					}

					grp.
						Default().
						Block(
							jen.Return(
								jenx.
									Errorf(
										fmt.Sprintf("%s: %%v is not a member of the enum", info.Name),
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
