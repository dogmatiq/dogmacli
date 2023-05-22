package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func (g *Generator) emitReifiedType(name string, t model.Type) {
	// g.pushScope(name)
	// model.VisitType(t, &reifyType{g, name})
	// g.popScope()
}

func (g *Generator) emitReifiedTypes() {
	unreified := g.unreified
	g.unreified = nil

	names := maps.Keys(unreified)
	slices.Sort(names)

	for _, n := range names {
		g.File.Line()
		g.emitReifiedType(n, unreified[n])
		g.emitReifiedTypes()
	}
}

func (g *Generator) nameFromScope() string {
	return normalize(g.scopes[len(g.scopes)-1]...)
}

func (g *Generator) reifyType(name string, t model.Type) {
	if g.reified == nil {
		g.reified = map[string]struct{}{}
	}

	if _, ok := g.reified[name]; !ok {
		g.reified[name] = struct{}{}

		if g.unreified == nil {
			g.unreified = map[string]model.Type{}
		}
		g.unreified[name] = t
	}
}

type reifyType struct {
	*Generator
	Name string
}

func (g *reifyType) Bool()                      { panic("not implemented") }
func (g *reifyType) Decimal()                   { panic("not implemented") }
func (g *reifyType) String()                    { panic("not implemented") }
func (g *reifyType) Integer()                   { panic("not implemented") }
func (g *reifyType) UInteger()                  { panic("not implemented") }
func (g *reifyType) DocumentURI()               { panic("not implemented") }
func (g *reifyType) URI()                       { panic("not implemented") }
func (g *reifyType) Null()                      { panic("not implemented") }
func (g *reifyType) Reference(*model.Reference) { panic("not implemented") }
func (g *reifyType) StringLit(*model.StringLit) { panic("not implemented") }
