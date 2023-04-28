package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) enum(
	code *jen.File,
	m metamodel.Enumeration,
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
		Add(g.typeRef(m.Type))

	code.
		Const().
		DefsFunc(func(code *jen.Group) {
			for _, member := range m.Members {
				documentation(code, member.Documentation)

				value := member.Value
				if v, ok := value.(float64); ok {
					value = int(v)
				}

				code.
					Id(normalizeName(m.Name) + normalizeName(member.Name)).
					Id(normalizeName(m.Name)).
					Op("=").
					Lit(value)
			}
		})
}
