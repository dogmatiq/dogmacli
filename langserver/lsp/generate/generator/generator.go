package generator

import (
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/generate/metamodel"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type generator struct {
	root          metamodel.Root
	names         map[string]typeInfo
	unionArities  map[int]struct{}
	tupleArities  map[int]struct{}
	namingContext []string
	pending       []jen.Code
}

type typeInfo struct {
	omittable bool
}

// Generate generates the Go representation of the LSP model.
func Generate(gen *jen.File) {
	g := &generator{
		root:  metamodel.Get(),
		names: map[string]typeInfo{},
	}

	for _, m := range g.root.Structures {
		g.names[normalizeName(m.Name)] = typeInfo{omittable: false}
	}
	for _, m := range g.root.Enumerations {
		g.names[normalizeName(m.Name)] = typeInfo{omittable: true}
	}
	for _, m := range g.root.TypeAliases {
		g.names[normalizeName(m.Name)] = typeInfo{omittable: g.isOmittable(m.Type)}
	}

	banner := func(text string) {
		gen.Comment(strings.Repeat("-", 72))
		gen.Comment(text)
		gen.Comment(strings.Repeat("-", 72))
		gen.Line()
	}

	banner("STRUCTURES")
	for _, m := range g.root.Structures {
		g.generateStruct(gen, m)
		g.flushPending(gen)
	}

	banner("ENUMERATIONS")
	for _, m := range g.root.Enumerations {
		g.generateEnum(gen, m)
		g.flushPending(gen)
	}

	banner("TYPE ALIASES")
	for _, m := range g.root.TypeAliases {
		g.generateAlias(gen, m)
		g.flushPending(gen)
	}

	if len(g.unionArities) != 0 {
		banner("UNIONS")
		g.generateUnions(gen)
	}

	if len(g.tupleArities) != 0 {
		banner("TUPLES")
		g.generateTuples(gen)
	}
}

func (g *generator) pushName(n string) {
	g.namingContext = append(g.namingContext, normalizeName(n))
}

func (g *generator) popName() {
	g.namingContext = g.namingContext[:len(g.namingContext)-1]
}

func (g *generator) uniqueName(desired string, info typeInfo) (actual, suffix string) {
	defer func() {
		g.names[actual] = info
	}()

	if _, ok := g.names[desired]; !ok {
		return desired, ""
	}

	i := 0
	for {
		i++

		suffix := strconv.Itoa(i)
		candidate := desired + suffix

		if _, ok := g.names[candidate]; !ok {
			return candidate, suffix
		}
	}
}

func (g *generator) flushPending(gen *jen.File) {
	for _, c := range g.pending {
		gen.Add(c)
	}
	g.pending = nil
}

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
	n = strings.Title(n)
	n = strings.ReplaceAll(n, "Uri", "URI")
	return n
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
