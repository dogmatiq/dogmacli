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
	embedded []*model.Struct,
	properties []model.Property,
) {
	g.emitStructType(name, embedded, properties)
	g.emitStructMarshalMethods(name, embedded, properties)
}

func (g *Generator) emitStructType(
	name string,
	embedded []*model.Struct,
	properties []model.Property,
) {
	g.File.
		Type().
		Id(name).
		StructFunc(func(grp *jen.Group) {
			for _, e := range embedded {
				grp.Id(identifier(e.TypeName))
			}

			if len(embedded) > 0 && len(properties) > 0 {
				grp.Line()
			}

			for _, p := range properties {
				if _, ok := p.Type.(model.StringLit); ok {
					// Don't add a struct field to represent string literal properties.
					// These are always handled at the (un)marshaling level.
					continue
				}

				g.enterProperty(p.Name)

				i := g.typeInfo(p.Type)
				t := g.typeExpr(p.Type)

				if p.Optional && i.UseOptional {
					t = jen.Id("Optional").Types(t)
				}

				documentation(grp, p.Documentation)

				grp.
					Id(identifier(p.Name)).
					Add(t)

				grp.Line()

				g.leaveProperty()
			}
		})
}

func (g *Generator) emitStructMarshalMethods(
	name string,
	embedded []*model.Struct,
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
			for _, e := range embedded {
				grp.
					If(
						jen.
							Err().
							Op(":=").
							Id("x").Dot(identifier(e.TypeName)).
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
				i := g.typeInfo(p.Type)

				fn := "marshalProperty"
				if p.Optional && i.UseOptional {
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
