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
	g.pending = append(
		g.pending,
		g.generateStructType(
			name,
			nil,
			t.LiteralStructProperties(),
		),
	)
}
