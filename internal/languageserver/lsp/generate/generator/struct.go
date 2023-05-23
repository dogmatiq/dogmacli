package generator

import (
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/jenx"
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

	g.emitStructUnmarshalMethods(
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
		Line().
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
						jenx.Errorf(
							fmt.Sprintf("%s: %%w", name),
							jen.Err(),
						),
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

	g.
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

func (g *generator) emitStructUnmarshalMethods(
	name string,
	embedded []model.Type,
	properties []*model.Property,
) {
	g.
		Line().
		Func().
		Params(
			jen.Id("x").Op("*").Id(name),
		).
		Id("UnmarshalJSON").
		Params(
			jen.Id("data").Index().Byte(),
		).
		Params(
			jen.Error(),
		).
		BlockFunc(func(g *jen.Group) {
			g.
				Var().
				Id("properties").
				Map(jen.String()).
				Qual("encoding/json", "RawMessage")

			g.
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
						jenx.Errorf(
							fmt.Sprintf("%s: %%w", name),
							jen.Err(),
						),
					),
				)

			g.
				If(
					jen.
						Err().
						Op(":=").
						Id("x").Dot("unmarshalProperties").
						Call(
							jen.Id("properties"),
						),
					jen.Err().Op("!=").Nil(),
				).
				Block(
					jen.Return(
						jenx.Errorf(
							fmt.Sprintf("%s: %%w", name),
							jen.Err(),
						),
					),
				)

			g.
				For().
				Id("k").Op(":=").Range().Id("properties").
				Block(
					jen.Return(
						jenx.Errorf(
							fmt.Sprintf("%s: %%s: unrecognized property", name),
							jen.Id("k"),
						),
					),
				)

			g.
				Return(
					jen.Nil(),
				)
		})

	g.
		Line().
		Func().
		Params(
			jen.Id("x").Op("*").Id(name),
		).
		Id("unmarshalProperties").
		Params(
			jen.Id("properties").Map(jen.String()).Qual("encoding/json", "RawMessage"),
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
							Dot("unmarshalProperties").
							Call(
								jen.Id("properties"),
							),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(jen.Err()),
					)
			}

			for _, p := range properties {
				g.
					IfFunc(func(g *jen.Group) {
						if t, ok := p.Type.(*model.StringLit); ok {
							g.
								Err().Op(":=").Id("unmarshalLiteralProperty").
								Call(
									jen.Id("properties"),
									jen.Lit(p.Name),
									jen.Lit(t.Value),
								)
						} else if kindOf(p) == reflect.Interface {
							fn := "unmarshalPropertyUsing"
							if p.Optional {
								fn = "unmarshalOptionalPropertyUsing"
							}

							g.
								Err().Op(":=").Id(fn).
								Call(
									jen.Id("properties"),
									jen.Lit(p.Name),
									jen.Op("&").Id("x").Dot(nameOf(p)),
									jen.Id("unmarshal"+nameOf(p.Type)),
								)
						} else {
							fn := "unmarshalProperty"
							if p.Optional {
								fn = "unmarshalOptionalProperty"
							}

							g.
								Err().Op(":=").Id(fn).
								Call(
									jen.Id("properties"),
									jen.Lit(p.Name),
									jen.Op("&").Id("x").Dot(nameOf(p)),
								)
						}

						g.Err().Op("!=").Nil()
					}).
					Block(
						jen.Return(jen.Err()),
					)
			}

			g.Return(jen.Nil())
		})
}
