package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) literalName(t *metamodel.Type) jen.Code {
	if len(g.namingContext) == 0 {
		panic("no context for naming literal")
	}

	key := strings.Join(g.namingContext, "")

	if g.pendingLiterals == nil {
		g.pendingLiterals = map[string][]*metamodel.Type{}
	}

	types := g.pendingLiterals[key]
	g.pendingLiterals[key] = append(types, t)

	return jen.Id(
		fmt.Sprintf("%s%d", key, len(types)),
	)
}

func (g *generator) literals(gen *jen.File) {
	for g.pendingLiterals != nil {
		literals := g.pendingLiterals
		g.pendingLiterals = nil

		for _, key := range sortedKeys(literals) {
			for i, t := range literals[key] {
				n := fmt.Sprintf("%s%d", key, i)

				g.pushName(n)

				gen.Type().
					Id(n).
					StructFunc(func(code *jen.Group) {
						for _, p := range t.LiteralProperties() {
							g.property(code, p)
						}
					})

				g.popName()
			}
		}
	}
}
