package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
	"golang.org/x/exp/slices"
)

// Generator generates Go code from the LSP meta-model.
type Generator struct {
	Model *model.Model
	File  *jen.File

	scopes    [][]string
	unreified map[string]model.Type
	reified   map[string]struct{}
}

// Generate populates g.File.
func (g *Generator) Generate() {
	defs := slices.Clone(g.Model.TypeDefs)
	slices.SortFunc(
		defs,
		func(a, b model.TypeDef) bool {
			return a.Name() < b.Name()
		},
	)

	for _, d := range defs {
		g.emitBanner(identifier(d.Name()))
		g.emitTypeDef(d)
		g.emitReifiedTypes()
	}

	methods := slices.Clone(g.Model.Methods)
	slices.SortFunc(
		methods,
		func(a, b model.Method) bool {
			return a.Name() < b.Name()
		},
	)

	for _, m := range methods {
		g.emitBanner(identifier(m.Name()))
		g.emitMethod(m)
		g.emitReifiedTypes()
	}
}

func (g *Generator) emitBanner(s string) {
	n := len(s)
	w := n + 8

	line := strings.Repeat("/", w)

	g.File.
		Line().
		Comment(line).
		Line().
		Commentf("/// %s ///", s).
		Line().
		Comment(line).
		Line()
}

func (g *Generator) pushScope(n string) {
	g.scopes = append(g.scopes, []string{n})
}

func (g *Generator) popScope() {
	g.scopes = g.scopes[:len(g.scopes)-1]
}

func (g *Generator) pushNestedScope(n string) {
	names := &g.scopes[len(g.scopes)-1]
	*names = append(*names, n)
}

func (g *Generator) popNestedScope() {
	names := &g.scopes[len(g.scopes)-1]
	*names = (*names)[:len(*names)-1]
}
