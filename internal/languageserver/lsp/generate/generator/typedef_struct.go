package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Struct(d model.Struct) {
	documentation(g.File, d.Documentation)
	g.emitStruct(
		identifier(d.TypeName),
		d.Embedded,
		d.Properties,
	)
}

func (g *Generator) emitStruct(
	name string,
	embedded []model.Type,
	properties []model.Property,
) {
	g.emitStructType(name, embedded, properties)
	g.emitStructMarshalMethods(name, embedded, properties)
}

func (g *Generator) emitStructType(
	name string,
	embedded []model.Type,
	properties []model.Property,
) {
	g.File.
		Type().
		Id(name).
		StructFunc(func(grp *jen.Group) {
			for _, t := range embedded {
				grp.Add(g.typeExpr(t))
			}

			if len(embedded) > 0 && len(properties) > 0 {
				grp.Line()
			}

			for _, p := range properties {
				g.pushNestedScope(p.Name)

				info := g.typeInfo(p.Type)

				if info.HasGoType {
					expr := g.typeExpr(p.Type)
					if p.Optional && !info.IsNativelyOptional {
						expr = jen.Id("Optional").Types(expr)
					}

					documentation(grp, p.Documentation)

					grp.
						Id(identifier(p.Name)).
						Add(expr)

					grp.Line()
				}

				g.popNestedScope()
			}
		})
}

func (g *Generator) emitStructMarshalMethods(
	name string,
	embedded []model.Type,
	properties []model.Property,
) {
	g.File.
		Func().
		Params(
			jen.Id("x").Id(name),
		).
		Id("MarshalJSON").
		Params().
		Params(
			jen.Index().Byte(),
			jen.Error(),
		).
		BlockFunc(func(grp *jen.Group) {
			grp.
				Var().
				Defs(
					jen.Id("w").Qual("bytes", "Buffer"),
					jen.Id("n").Int(),
				)

			grp.
				Line().
				Id("w").
				Dot("WriteByte").
				Call(
					jen.LitRune('{'),
				)

			grp.
				Line().
				If(
					jen.
						Err().
						Op(":=").
						Id("x").Dot("marshalProperties").
						Call(
							jen.Op("&").Id("w"),
							jen.Op("&").Id("n"),
						),
					jen.Err().Op("!=").Nil(),
				).
				Block(
					jen.Return(
						jen.Nil(),
						jen.Err(),
					),
				)

			grp.
				Line().
				Id("w").
				Dot("WriteByte").
				Call(
					jen.LitRune('}'),
				)

			grp.
				Line().
				Return(
					jen.Id("w").Dot("Bytes").Call(),
					jen.Nil(),
				)
		})

	g.File.
		Line().
		Func().
		Params(
			jen.Id("x").Id(name),
		).
		Id("marshalProperties").
		Params(
			jen.Id("w").Op("*").Qual("bytes", "Buffer"),
			jen.Id("n").Op("*").Int(),
		).
		Params(
			jen.Error(),
		).
		BlockFunc(func(grp *jen.Group) {
			for _, t := range embedded {
				grp.
					If(
						jen.
							Err().
							Op(":=").
							Id("x").Op(".").Add(g.typeExpr(t)).
							Dot("marshalProperties").
							Call(
								jen.Id("w"),
								jen.Id("n"),
							),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(jen.Err()),
					)
			}

			for _, p := range properties {
				info := g.typeInfo(p.Type)

				fn := "marshalProperty"
				if p.Optional && !info.IsNativelyOptional {
					fn = "marshalOptionalProperty"
				}

				expr := jen.Id("x").Dot(identifier(p.Name))
				if t, ok := p.Type.(model.StringLit); ok {
					expr = jen.Lit(t.Value)
				}

				grp.
					If(
						jen.Err().Op(":=").Id(fn).Call(
							jen.Id("w"),
							jen.Id("n"),
							jen.Lit(p.Name),
							expr,
						),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(jen.Err()),
					)
			}

			grp.Return(jen.Nil())
		})
}
