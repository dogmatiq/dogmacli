package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) alias(
	code *jen.File,
	m metamodel.TypeAlias,
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
}
