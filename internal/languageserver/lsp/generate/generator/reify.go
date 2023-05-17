package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func (g *Generator) emitReifiedType(n string, t model.Type) {
	g.enterType(n)
	model.VisitType(t, &reifyType{g, n})
	g.leaveType()
}

func (g *Generator) emitReifiedTypes() {
	literals := g.literals
	g.literals = nil

	names := maps.Keys(literals)
	slices.Sort(names)

	for _, n := range names {
		g.File.Line()
		g.emitReifiedType(n, literals[n])
		g.emitReifiedTypes()
	}
}

func (g *Generator) enterType(n string) {
	g.scopes = append(g.scopes, []string{n})
}

func (g *Generator) leaveType() {
	g.scopes = g.scopes[:len(g.scopes)-1]
}

func (g *Generator) enterProperty(n string) {
	names := &g.scopes[len(g.scopes)-1]
	*names = append(*names, n)
}

func (g *Generator) leaveProperty() {
	names := &g.scopes[len(g.scopes)-1]
	*names = (*names)[:len(*names)-1]
}

func (g *Generator) reify(t model.Type) *jen.Statement {
	name := identifier(g.scopes[len(g.scopes)-1]...)

	if g.literals == nil {
		g.literals = map[string]model.Type{}
	}

	g.literals[name] = t

	return jen.Id(name)
}

type reifyType struct {
	*Generator
	Name string
}

func (g *reifyType) Bool()                       { panic("not implemented") }
func (g *reifyType) Decimal()                    { panic("not implemented") }
func (g *reifyType) String()                     { panic("not implemented") }
func (g *reifyType) Integer()                    { panic("not implemented") }
func (g *reifyType) UInteger()                   { panic("not implemented") }
func (g *reifyType) DocumentURI()                { panic("not implemented") }
func (g *reifyType) URI()                        { panic("not implemented") }
func (g *reifyType) Null()                       { panic("not implemented") }
func (g *reifyType) Reference(model.Reference)   { panic("not implemented") }
func (g *reifyType) Array(model.Array)           { panic("not implemented") }
func (g *reifyType) Map(model.Map)               { panic("not implemented") }
func (g *reifyType) StringLit(t model.StringLit) { panic("not implemented") }
