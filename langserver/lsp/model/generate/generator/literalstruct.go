package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) literalStructRef(t *metamodel.Type) jen.Code {
	if len(g.namingContext) == 0 {
		panic("no context for naming literal")
	}

	name, suffix := g.uniqueName(
		strings.Join(g.namingContext, ""),
	)

	g.pushName(suffix)
	defer g.popName()

	// Add the statement to the pending list _before_ generating its body, which
	// may contain further literals.
	gen := &jen.Statement{}
	g.pending = append(g.pending, gen)

	gen.Type().
		Id(name).
		StructFunc(func(gen *jen.Group) {
			for _, p := range t.LiteralStructProperties() {
				g.generateStructProperty(gen, p)
			}
		})

	return jen.Id(name)
}
