package main

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateStructs(gen *jen.File) {
	generateBanner(gen, "STRUCTURES")

	for _, m := range g.root.Structures {
		g.generateStruct(gen, m)
		g.flushPending(gen)
	}
}

func (g *generator) generateStruct(
	gen *jen.File,
	m metamodel.Structure,
) {
	g.pushName(m.Name)
	defer g.popName()

	if m.Documentation == "" {
		gen.Line()
	} else {
		generateDocs(gen, m.Documentation)
	}

	name := normalizeName(m.Name)

	gen.Add(
		g.generateStructType(
			name,
			m.Embeds(),
			m.Properties,
		),
	)
}

func (g *generator) generateStructType(
	name string,
	embeds []*metamodel.Type,
	properties []metamodel.Property,
) jen.Code {
	gen := &jen.Statement{}

	gen.Type().
		Id(name).
		StructFunc(func(gen *jen.Group) {
			for _, p := range embeds {
				gen.Id(normalizeName(p.Name))
			}
			for _, p := range properties {
				g.generateStructProperty(gen, p)
			}
		})

	gen.Line().
		Func().
		Params(jen.Id("x").Op("*").Id(name)).
		Id("UnmarshalJSON").
		Params(
			jen.Id("data").Index().Byte(),
		).
		Params(
			jen.Error(),
		).
		BlockFunc(func(gen *jen.Group) {
			gen.Type().Id("plain").Id(name)
			gen.If(
				jen.Err().Op(":=").Id("unmarshal").Call(
					jen.Id("data"),
					jen.Parens(
						jen.Op("*").Id("plain"),
					).Call(
						jen.Id("x"),
					),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)

			for _, p := range properties {
				if p.Optional {
					continue
				}

				gen.Line()

				if zero, ok := g.zeroValue(p.Type); ok {
					gen.If(
						jen.
							Id("x").Dot(normalizeName(p.Name)).
							Op("==").
							Add(zero),
					).Block(
						jen.Return(
							jen.Qual("errors", "New").Call(
								jen.Lit(
									fmt.Sprintf(
										"%q property is required",
										p.Name,
									),
								),
							),
						),
					)
				} else {
					gen.Var().Id("_").Qual("encoding/json", "Unmarshaler").
						Op("=").
						Op("&").Id("x").Dot(normalizeName(p.Name))
				}
			}

			gen.Line().
				Return(jen.Nil())
		})

	return gen
}

func (g *generator) generateStructProperty(
	gen *jen.Group,
	m metamodel.Property,
) {
	g.pushName(m.Name)
	defer g.popName()

	generateDocs(gen, m.Documentation)

	expr := g.typeExpr(m.Type)
	tag := m.Name

	if m.Optional {
		tag += ",omitempty"

		if !g.isOmittable(m.Type) {
			expr = jen.
				Op("*").
				Add(expr)
		}
	}

	gen.Id(normalizeName(m.Name)).
		Add(expr).
		Tag(map[string]string{
			"json": tag,
		})
}
