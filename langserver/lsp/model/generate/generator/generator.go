package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type generator struct {
	root            metamodel.Root
	unionArities    map[int]struct{}
	tupleArities    map[int]struct{}
	pendingLiterals map[string][]*metamodel.Type
	namingContext   []string
}

// Generate generates the Go representation of the LSP model.
func Generate(code *jen.File) {
	g := &generator{
		root: metamodel.Get(),
	}

	banner := func(text string) {
		code.Comment(strings.Repeat("-", 72))
		code.Comment(text)
		code.Comment(strings.Repeat("-", 72))
		code.Line()
	}

	banner("STRUCTURES")
	for _, m := range g.root.Structures {
		g.structure(code, m)
		g.literals(code)
	}

	banner("ENUMERATIONS")
	for _, m := range g.root.Enumerations {
		g.enum(code, m)
		g.literals(code)
	}

	banner("TYPE ALIASES")
	for _, m := range g.root.TypeAliases {
		g.alias(code, m)
		g.literals(code)
	}

	banner("UNIONS & TUPLES")
	g.unions(code)
	g.tuples(code)
}

func documentation(
	code interface {
		Comment(string) *jen.Statement
	},
	docs string,
) {
	if docs != "" {
		for _, line := range strings.Split(docs, "\n") {
			code.Comment(line)
		}
	}
}

func (g *generator) hasSymbol(n string) bool {
	for _, def := range g.root.Structures {
		if def.Name == n {
			return true
		}
	}

	for _, def := range g.root.Enumerations {
		if def.Name == n {
			return true
		}
	}

	for _, def := range g.root.TypeAliases {
		if def.Name == n {
			return true
		}
	}

	_, ok := g.pendingLiterals[n]

	return ok
}

func normalizeName(n string) string {
	n = strings.Title(n)
	n = strings.ReplaceAll(n, "Uri", "URI")
	return n
}

func (g *generator) pushName(n string) {
	g.namingContext = append(g.namingContext, normalizeName(n))
}

func (g *generator) popName() {
	g.namingContext = g.namingContext[:len(g.namingContext)-1]
}

func sortedKeys[M map[K]V, K constraints.Ordered, V any](m M) []K {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}
