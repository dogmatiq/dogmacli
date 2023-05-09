package main

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) literalStructTypeExpr(t *metamodel.Type) jen.Code {
	if len(g.namingContext) == 0 {
		panic("no context for naming literal")
	}

	name := strings.Join(g.namingContext, "")
	name, suffix := g.generateUniqueName(name)

	g.pushName(suffix)
	defer g.popName()

	g.generateLiteralStruct(name, t)

	return jen.Id(name)
}

func (g *generator) generateLiteralStruct(name string, t *metamodel.Type) {
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
}
