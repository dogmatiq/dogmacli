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
		Params(jen.Id("x").Id(name)).
		Id("Validate").
		Params().
		Params(
			jen.Error(),
		).
		BlockFunc(func(gen *jen.Group) {
			for _, p := range embeds {
				gen.
					If(
						jen.Err().
							Op(":=").
							Id("x").Dot(normalizeName(p.Name)).
							Dot("Validate").Call(),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(jen.Err()),
					)
			}

			for _, p := range properties {
				pname := normalizeName(p.Name)
				info := g.typeInfo(p.Type)

				validate := jen.
					If(
						jen.Err().
							Op(":=").
							Id("x").Dot(pname).
							Dot("Validate").Call(),
						jen.Err().Op("!=").Nil(),
					).
					Block(
						jen.Return(
							jen.Qual("fmt", "Errorf").Call(
								jen.Lit(
									fmt.Sprintf("invalid field %q: %%w", pname),
								),
								jen.Err(),
							),
						),
					)

				if p.Optional {
					if info.IsValidateable {
						if info.AddPointer {
							gen.
								If(
									jen.Id("x").Dot(pname).
										Op("!=").
										Nil(),
								).
								Block(validate)
						} else {
							gen.Add(validate)
						}
					}
				} else {
					if info.AddPointer {
						gen.
							If(
								jen.Id("x").Dot(pname).
									Op("==").
									Nil(),
							).
							Block(
								jen.Return(
									jen.Qual("errors", "New").Call(
										jen.Lit(
											fmt.Sprintf(
												"missing required field %q",
												p.Name,
											),
										),
									),
								),
							)
					}

					if info.IsValidateable {
						gen.Add(validate)
					}
				}
			}

			gen.Return(jen.Nil())
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

	tag := m.Name
	if m.Optional {
		tag += ",omitempty"
	}

	expr := g.typeExpr(m.Type)
	info := g.typeInfo(m.Type)

	if info.AddPointer {
		expr = jen.Op("*").Add(expr)
	}

	gen.Id(normalizeName(m.Name)).
		Add(expr).
		Tag(map[string]string{
			"json": tag,
		})
}
