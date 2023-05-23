package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Generate produces Go representations of the given LSP model.
func Generate(
	m *model.Model,
	f *jen.File,
) {
	g := &generator{
		Group: f.Group,
		Model: m,
	}

	model.VisitNode(m, g)
}

type generator struct {
	*jen.Group
	Model *model.Model
}

// withGroup returns a function that can be passed to jen's <Elem>Func() methods
// which, when invoked calls fn() with g.Group set to the provided by the
// XXXFunc() method.
func (g *generator) withGroup(
	fn func(),
) func(*jen.Group) {
	return func(after *jen.Group) {
		before := g.Group
		g.Group = after
		fn()
		g.Group = before
	}
}

func (g *generator) VisitModel(n *model.Model) {
	types := map[string]model.Type{}
	for _, t := range n.Types {
		if _, ok := t.(*model.Reference); ok {
			continue
		}

		if name, ok := tryNameOf(t); ok {
			types[name] = t
		}
	}

	names := append(
		maps.Keys(types),
		maps.Keys(n.Defs)...,
	)
	slices.Sort(names)
	names = slices.Compact(names)

	for _, name := range names {
		if d, ok := n.Defs[name]; ok {
			model.VisitNode(d, g)
		}
		if t, ok := types[name]; ok {
			model.VisitNode(t, g)
		}
	}
}

func (g *generator) VisitCall(n *model.Call)                 {}
func (g *generator) VisitNotification(n *model.Notification) {}
func (g *generator) VisitReference(n *model.Reference)       {}
func (g *generator) VisitStringLit(n *model.StringLit)       {}
