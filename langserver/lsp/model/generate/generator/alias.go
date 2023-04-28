package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) generateAlias(
	gen *jen.File,
	m metamodel.TypeAlias,
) {
	if strings.HasPrefix(m.Name, "LSP") {
		return
	}

	g.pushName(m.Name)
	defer g.popName()

	if m.Documentation == "" {
		gen.Line()
	} else {
		generateDocs(gen, m.Documentation)
	}

	gen.Type().
		Id(normalizeName(m.Name)).
		Add(g.typeRef(m.Type))
}
