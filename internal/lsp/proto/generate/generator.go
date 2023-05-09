package main

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type generator struct {
	root          metamodel.Root
	names         map[string]struct{}
	unionArities  map[int]struct{}
	tupleArities  map[int]struct{}
	namingContext []string
	pending       []jen.Code
}

// generate generates the Go representation of the LSP model.
func generate(gen *jen.File) {
	g := &generator{
		root:         metamodel.Get(),
		names:        map[string]struct{}{},
		unionArities: map[int]struct{}{},
		tupleArities: map[int]struct{}{},
	}

	for _, m := range g.root.Structures {
		g.names[normalizeName(m.Name)] = struct{}{}
	}
	for _, m := range g.root.Enumerations {
		g.names[normalizeName(m.Name)] = struct{}{}
	}
	for _, m := range g.root.TypeAliases {
		g.names[normalizeName(m.Name)] = struct{}{}
	}

	g.generateRequests(gen)
	g.generateNotifications(gen)

	g.generateStructs(gen)
	g.generateEnums(gen)
	g.generateAliases(gen)

	g.generateUnions(gen)
	g.generateTuples(gen)
}

func (g *generator) pushName(name string) {
	g.namingContext = append(g.namingContext, normalizeName(name))
}

func (g *generator) popName() {
	g.namingContext = g.namingContext[:len(g.namingContext)-1]
}

func (g *generator) flushPending(gen *jen.File) {
	for _, c := range g.pending {
		gen.Add(c)
	}
	g.pending = nil
}

// generateBanner writes a generateBanner comment to the generated file.
func generateBanner(gen *jen.File, text string) {
	sep := strings.Repeat("-", 120-3)
	gen.Comment(sep)
	gen.Comment(text)
	gen.Comment(sep)
	gen.Line()
}

// generateDocs writes documentation comments to the generated file.
func generateDocs(
	gen interface{ Comment(string) *jen.Statement },
	docs string,
) {
	if docs != "" {
		for _, line := range strings.Split(docs, "\n") {
			gen.Comment(line)
		}
	}
}

func normalizeName(n string) string {
	n = strings.ReplaceAll(n, "/", " ")
	n = strings.ReplaceAll(n, "$", " ")
	n = strings.Title(n)
	n = strings.ReplaceAll(n, " ", "")
	n = strings.ReplaceAll(n, "Uri", "URI")
	return n
}

func normalizeUnexportedName(n string) string {
	n = normalizeName(n)
	return strings.ToLower(n[:1]) + n[1:]
}

func sortedKeys[M map[K]V, K constraints.Ordered, V any](m M) []K {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}

var ordinalFieldNames = []string{
	"First",
	"Second",
	"Third",
	"Fourth",
	"Fifth",
	"Sixth",
	"Seventh",
	"Eighth",
	"Ninth",
	"Tenth",
}
