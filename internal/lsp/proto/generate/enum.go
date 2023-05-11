package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateEnums(gen *jen.File) {
	generateBanner(gen, "ENUMERATIONS")

	for _, m := range g.root.Enumerations {
		g.generateEnum(gen, m)
		g.flushPending(gen)
	}
}

func (g *generator) generateEnum(
	gen *jen.File,
	m metamodel.Enumeration,
) {
	g.pushName(m.Name)
	defer g.popName()

	if m.Documentation == "" {
		gen.Line()
	} else {
		generateDocs(gen, m.Documentation)
	}

	name := normalizeName(m.Name)

	gen.Type().
		Id(name).
		Add(g.typeExpr(m.Type))

	gen.Const().
		DefsFunc(func(gen *jen.Group) {
			for _, member := range m.Members {
				generateDocs(gen, member.Documentation)

				value := member.Value
				if v, ok := value.(float64); ok {
					value = int(v)
				}

				gen.Id(name + normalizeName(member.Name)).
					Id(name).
					Op("=").
					Lit(value)
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
		Block(
			jen.Switch(jen.Id("x")).
				BlockFunc(func(gen *jen.Group) {
					for _, member := range m.Members {
						gen.Case(
							jen.Id(name + normalizeName(member.Name)),
						).Block(
							jen.Return(jen.Nil()),
						)
					}

					gen.Default().
						Block(
							jen.Return(
								jen.Qual("fmt", "Errorf").
									Call(
										jen.Lit("%#v is not a valid "+name),
										jen.Id("x"),
									),
							),
						)
				}),
		)
}
