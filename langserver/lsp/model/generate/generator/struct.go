package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) structure(
	code *jen.File,
	m metamodel.Structure,
) {
	g.pushName(m.Name)
	defer g.popName()

	if m.Documentation == "" {
		code.Line()
	} else {
		documentation(code, m.Documentation)
	}

	code.
		Type().
		Id(normalizeName(m.Name)).
		StructFunc(func(code *jen.Group) {
			for _, p := range m.Extends {
				code.Id(normalizeName(p.Name))
			}
			for _, p := range m.Mixins {
				code.Id(normalizeName(p.Name))
			}
			for _, p := range m.Properties {
				g.property(code, p)
			}
		})
}

func (g *generator) property(
	code *jen.Group,
	m metamodel.Property,
) {
	g.pushName(m.Name)
	defer g.popName()

	documentation(code, m.Documentation)

	ref := g.typeRef(m.Type)
	tag := m.Name

	if m.Optional {
		tag += ",omitempty"

		switch m.Type.Kind {
		case "base":
		case "map":
		case "array":
		default:
			ref = jen.
				Op("*").
				Add(ref)
		}
	}

	code.
		Id(normalizeName(m.Name)).
		Add(ref).
		Tag(map[string]string{
			"json": tag,
		})
}

var ordinalFieldNames = []string{
	"First",
	"Second",
	"Third",
	"Fourth",
	"Fifth",
	"Sixth",
	"Seventh",
	"Eighth",
	"Ninth",
	"Tenth",
}
