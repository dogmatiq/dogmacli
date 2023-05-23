package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *generator) VisitStruct(n *model.Struct) {
	name := nameOf(n)

	documentation(
		g,
		n.Documentation(),
		"%s is a named structure definition.",
		name,
	)

	g.emitStruct(
		name,
		n.EmbeddedTypes,
		n.Properties,
	)
}

func (g *generator) emitStruct(
	name string,
	embedded []model.Type,
	properties []*model.Property,
) {
	g.
		Type().
		Id(name).
		StructFunc(
			g.withGroup(func() {
				for _, t := range embedded {
					g.Id(nameOf(t))
				}

				if len(embedded) > 0 && len(properties) > 0 {
					g.Line()
				}

				for _, p := range properties {
					model.VisitNode(p, g)
				}
			}),
		)

	g.emitStructMarshalMethods(
		name,
		embedded,
		properties,
	)
}

func (g *generator) VisitProperty(n *model.Property) {
	name := nameOf(n)
	typ, ok := tryNameOf(n.Type)
	if !ok {
		return
	}

	documentation(
		g,
		n.Documentation,
		"",
	)

	expr := jen.Id(typ)
	if n.Optional {
		expr = jen.Id("Optional").Types(expr)
	}

	g.
		Id(name).
		Add(expr)
}

func (g *generator) emitStructMarshalMethods(
	name string,
	embedded []model.Type,
	properties []*model.Property,
) {
	g.
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
		BlockFunc(func(g *jen.Group) {
			g.
				Var().
				Defs(
					jen.Id("w").Qual("bytes", "Buffer"),
					jen.Id("n").Int(),
				)

			g.
				Line().
				Id("w").Dot("WriteByte").
				Call(jen.LitRune('{'))

			g.
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

			g.
				Id("w").Dot("WriteByte").
				Call(jen.LitRune('}'))

			g.
				Line().
				Return(
					jen.Id("w").Dot("Bytes").Call(),
					jen.Nil(),
				)
		})

	g.Line()

	g.
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
		BlockFunc(func(g *jen.Group) {
			for _, t := range embedded {
				g.
					If(
						jen.
							Err().
							Op(":=").
							Id("x").Op(".").Id(nameOf(t)).
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
				fn := "marshalProperty"
				if p.Optional {
					fn = "marshalOptionalProperty"
				}

				expr := jen.Id("x").Dot(nameOf(p))
				if t, ok := p.Type.(*model.StringLit); ok {
					expr = jen.Lit(t.Value)
				}

				g.
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

			g.Return(jen.Nil())
		})
}
