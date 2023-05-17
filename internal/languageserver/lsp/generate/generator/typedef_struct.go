package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/jenx"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Struct(d model.Struct) {
	documentation(g.File, d.Documentation)
	g.File.
		Type().
		Id(identifier(d.TypeName)).
		StructFunc(func(grp *jen.Group) {
			for _, e := range d.Embedded {
				grp.Id(identifier(e.TypeName))
			}

			if len(d.Embedded) > 0 && len(d.Properties) > 0 {
				grp.Line()
			}

			for _, p := range d.Properties {
				g.structProperty(grp, p)
				grp.Line()
			}
		})

	g.File.Line()
	g.structMarshalMethods(d)

	// g.File.Line()
	// g.structUnmarshalMethod(d)
}

func (g *Generator) structProperty(grp *jen.Group, p model.Property) {
	if _, ok := p.Type.(model.StringLit); ok {
		return
	}

	g.pushName(p.Name)
	defer g.popName()

	documentation(grp, p.Documentation)

	i := g.typeInfo(p.Type)
	t := g.typeExpr(p.Type)

	if p.Optional && i.UseOptional {
		t = jen.Id("Optional").Types(t)
	}

	grp.
		Id(identifier(p.Name)).
		Add(t)
}

func (g *Generator) structMarshalMethods(d model.Struct) {
	g.File.
		Func().
		Params(
			jen.Id("x").Id(identifier(d.TypeName)),
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
			jen.Id("x").Id(identifier(d.TypeName)),
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
			for _, e := range d.Embedded {
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

			for _, p := range d.Properties {
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

func (g *Generator) structUnmarshalMethod(d model.Struct) {
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
				Var().
				Id("properties").
				Map(jen.String()).
				Qual("encoding/json", "RawMessage")

			grp.
				If(
					jen.
						Err().
						Op(":=").
						Qual("encoding/json", "Unmarshal").
						Call(
							jen.Id("data"),
							jen.Op("&").Id("properties"),
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

			// for _, e := range d.Embedded {
			// }

			// for _, p := range d.Properties {
			// 	grp.Comment(p.Name)
			// }

			grp.
				Line().
				Return(jen.Nil())
		})
}
