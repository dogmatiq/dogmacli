package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
	"golang.org/x/exp/slices"
)

// Generator generates Go code from the LSP meta-model.
type Generator struct {
	Model model.Model
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
