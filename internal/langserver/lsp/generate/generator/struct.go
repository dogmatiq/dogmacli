package generator

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/jenx"
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
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
	if n.IsOptional && kindOf(n) == reflect.Struct {
		expr = jen.Op("*").Id(typ)
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
				g.
					If(
						jen.
							Err().Op(":=").Id("x").Dot("marshal"+nameOf(p)+"Property").
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

			g.Return(jen.Nil())
		})

	for _, p := range properties {
		g.
			Line().
			Func().
			Params(
				jen.Id("x").Id(name),
			).
			Id("marshal"+nameOf(p)+"Property").
			Params(
				jen.Id("w").Op("*").Qual("bytes", "Buffer"),
				jen.Id("n").Op("*").Int(),
			).
			Params(
				jen.Error(),
			).
			BlockFunc(func(g *jen.Group) {
				g.
					If(
						jen.Op("*").Id("n").Op("++"),
						jen.Op("*").Id("n").Op(">").Lit(1),
					).
					Block(
						jen.
							Id("w").Dot("WriteByte").
							Call(jen.LitRune(',')),
					)

				key, err := json.Marshal(p.Name)
				if err != nil {
					panic(err)
				}

				if t, ok := p.Type.(*model.StringLit); ok {
					value, err := json.Marshal(t.Value)
					if err != nil {
						panic(err)
					}

					g.
						Id("w").Dot("WriteString").
						Call(jen.Lit(string(key) + ":" + string(value)))

					g.Return(jen.Nil())
					return
				}

				prop := jen.Id("x").Dot(nameOf(p))

				if p.IsOptional {
					g.
						IfFunc(func(g *jen.Group) {
							switch kindOf(p) {
							case reflect.Slice, reflect.String, reflect.Map:
								g.Len(prop).Op("==").Lit(0)
							case reflect.Bool:
								g.Op("!").Add(prop)
							case reflect.Int32, reflect.Uint32, reflect.Float64:
								g.Add(prop).Op("==").Lit(0)
							default:
								g.Add(prop).Op("==").Nil()
							}
						}).
						Block(
							jen.Return(jen.Nil()),
						)
				}

				g.
					Id("data").Op(",").Err().
					Op(":=").
					Qual("encoding/json", "Marshal").
					Call(prop)

				g.
					If(
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(jen.Err()),
					)

				g.
					Id("w").Dot("WriteString").
					Call(jen.Lit(string(key) + ":"))

				g.
					Id("w").Dot("Write").
					Call(jen.Id("data"))

				g.Return(jen.Nil())
			})
	}
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
				Id("p").
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
							jen.Op("&").Id("p"),
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
							jen.Id("p"),
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
				Id("k").Op(":=").Range().Id("p").
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
			jen.Id("p").Map(jen.String()).Qual("encoding/json", "RawMessage"),
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
								jen.Id("p"),
							),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(jen.Err()),
					)
			}

			for _, p := range properties {
				g.
					If(
						jen.
							Err().Op(":=").Id("x").Dot("unmarshal"+nameOf(p)+"Property").
							Call(
								jen.Id("p"),
							),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(
							jenx.Errorf(
								fmt.Sprintf("%s: %%w", p.Name),
								jen.Err(),
							),
						),
					)
			}

			g.Return(jen.Nil())
		})

	for _, p := range properties {
		g.
			Line().
			Func().
			Params(
				jen.Id("x").Id(name),
			).
			Id("unmarshal" + nameOf(p) + "Property").
			Params(
				jen.Id("p").Map(jen.String()).Qual("encoding/json", "RawMessage"),
			).
			Params(
				jen.Error(),
			).
			BlockFunc(func(g *jen.Group) {
				g.
					If(
						jen.
							Id("data").Op(",").Id("ok").
							Op(":=").
							Id("p").Index(jen.Lit(p.Name)),
						jen.Id("ok"),
					).
					BlockFunc(func(g *jen.Group) {
						if t, ok := p.Type.(*model.StringLit); ok {
							g.
								Var().
								Id("v").
								String()

							g.
								If(
									jen.
										Err().
										Op(":=").
										Qual("encoding/json", "Unmarshal").
										Call(
											jen.Id("data"),
											jen.Op("&").Id("v"),
										),
									jen.Err().Op("!=").Nil(),
								).
								Block(
									jen.Return(jen.Err()),
								)

							g.If(
								jen.Id("v").Op("!=").Lit(t.Value),
							).Block(
								jen.Return(
									jenx.Errorf(
										fmt.Sprintf("unexpected value %%q, expected %q", t.Value),
										jen.Id("v"),
									),
								),
							)

							g.Return(jen.Nil())
							return
						}

						fn := jen.Qual("encoding/json", "Unmarshal")
						if kindOf(p) == reflect.Interface {
							fn = jen.Id("unmarshal" + nameOf(p.Type))
						}

						g.Return(
							jen.
								Add(fn).
								Call(
									jen.Id("data"),
									jen.Op("&").Id("x").Dot(nameOf(p)),
								),
						)
					})

				if p.IsOptional {
					g.Return(jen.Nil())
				} else {
					g.Return(
						jen.
							Qual("errors", "New").
							Call(jen.Lit("mandatory property is not present")),
					)
				}
			})
	}
}

// for _, p := range properties {
// 	g.
// 		IfFunc(func(g *jen.Group) {
// 			if t, ok := p.Type.(*model.StringLit); ok {
// 				g.
// 					Err().Op(":=").Id("unmarshalLiteralProperty").
// 					Call(
// 						jen.Id("properties"),
// 						jen.Lit(p.Name),
// 						jen.Lit(t.Value),
// 					)
// 			} else if kindOf(p) == reflect.Interface {
// 				fn := "unmarshalMandatoryPropertyUsing"
// 				if p.IsOptional {
// 					if kindOf(p) == reflect.Struct {
// 						fn = "unmarshalPointerToOptionalPropertyUsing"
// 					} else {
// 						fn = "unmarshalOptionalPropertyUsing"
// 					}
// 				}

// 				g.
// 					Err().Op(":=").Id(fn).
// 					Call(
// 						jen.Id("properties"),
// 						jen.Lit(p.Name),
// 						jen.Op("&").Id("x").Dot(nameOf(p)),
// 						jen.Id("unmarshal"+nameOf(p.Type)),
// 					)
// 			} else {
// 				fn := "unmarshalMandatoryProperty"
// 				if p.IsOptional {
// 					if kindOf(p) == reflect.Struct {
// 						fn = "unmarshalPointerToOptionalProperty"
// 					} else {
// 						fn = "unmarshalOptionalProperty"
// 					}
// 				}

// 				g.
// 					Err().Op(":=").Id(fn).
// 					Call(
// 						jen.Id("properties"),
// 						jen.Lit(p.Name),
// 						jen.Op("&").Id("x").Dot(nameOf(p)),
// 					)
// 			}

// 			g.Err().Op("!=").Nil()
// 		}).
// 		Block(
// 			jen.Return(jen.Err()),
// 		)
// }
