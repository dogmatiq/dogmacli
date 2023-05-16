package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func (g *Generator) typeLit(n string, t model.Type) {
	g.pushScope(n)
	defer g.popScope()

	model.VisitType(t, &typeLit{g, n})
}

func (g *Generator) pushScope(n string) {
	g.scopes = append(g.scopes, []string{n})
}

func (g *Generator) popScope() {
	g.scopes = g.scopes[:len(g.scopes)-1]
}

func (g *Generator) pushName(n string) {
	names := &g.scopes[len(g.scopes)-1]
	*names = append(*names, n)
}

func (g *Generator) popName() {
	names := &g.scopes[len(g.scopes)-1]
	*names = (*names)[:len(*names)-1]
}

func (g *Generator) enqueueLiteral(t model.Type) string {
	name := identifier(g.scopes[len(g.scopes)-1]...)

	if g.literals == nil {
		g.literals = map[string]model.Type{}
	}

	g.literals[name] = t

	return name
}

func (g *Generator) flushLiterals() {
	literals := g.literals
	g.literals = nil

	names := maps.Keys(literals)
	slices.Sort(names)

	for _, n := range names {
		g.File.Line()
		g.typeLit(n, literals[n])
		g.flushLiterals()
	}
}

type typeLit struct {
	*Generator
	Name string
}

func (g *typeLit) Bool()                       { panic("not implemented") }
func (g *typeLit) Decimal()                    { panic("not implemented") }
func (g *typeLit) String()                     { panic("not implemented") }
func (g *typeLit) Integer()                    { panic("not implemented") }
func (g *typeLit) UInteger()                   { panic("not implemented") }
func (g *typeLit) DocumentURI()                { panic("not implemented") }
func (g *typeLit) URI()                        { panic("not implemented") }
func (g *typeLit) Null()                       { panic("not implemented") }
func (g *typeLit) Reference(model.Reference)   { panic("not implemented") }
func (g *typeLit) Array(model.Array)           { panic("not implemented") }
func (g *typeLit) Map(model.Map)               { panic("not implemented") }
func (g *typeLit) StringLit(t model.StringLit) { panic("not implemented") }
