package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

// Generator generates Go code from the LSP meta-model.
type Generator struct {
	Model *model.Model
	File  *jen.File

	scopes    [][]string
	unreified map[string]model.Type
	reified   map[string]struct{}
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
