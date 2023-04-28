package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) literalName(t *metamodel.Type) jen.Code {
	if len(g.namingContext) == 0 {
		panic("no context for naming literal")
	}

	name := "Lit" + strings.Join(g.namingContext, "")

	if g.pendingLiterals == nil {
		g.pendingLiterals = map[string]*metamodel.Type{}
	}

	g.pendingLiterals[name] = t

	return jen.Id(name)
}

func (g *generator) literals(code *jen.File) {
	for g.pendingLiterals != nil {
		literals := g.pendingLiterals
		g.pendingLiterals = nil

		for _, n := range sortedKeys(literals) {
			t := literals[n]

			code.
				Type().
				Id(n).
				StructFunc(func(code *jen.Group) {
					for _, p := range t.LiteralProperties() {
						g.property(code, p)
					}
				})
		}
	}
}
