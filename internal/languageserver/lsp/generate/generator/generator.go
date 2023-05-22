package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
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
	}

	model.VisitNode(m, g)
}

type generator struct {
	*jen.Group
}

// withGroup returns a function that can be passed to jen's XXXFunc() methods
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

func (g *generator) VisitArray(n *model.Array) {
	name := nameOf(n)
	elem := nameOf(n.ElementType)

	g.
		Commentf("%s is an array of %s elements.", name, elem).
		Line().
		Type().
		Id(name).
		Index().
		Id(elem)
}

func (g *generator) VisitMap(n *model.Map) {
	name := nameOf(n)
	key := nameOf(n.KeyType)
	value := nameOf(n.ValueType)

	g.
		Commentf("%s is an array of %s to %s.", name, key, value).
		Line().
		Type().
		Id(name).
		Map(jen.Id(key)).
		Id(value)
}

func (g *generator) VisitAnd(n *model.And) {
	name := nameOf(n)

	g.
		Commentf("%s is the intersection of several types.", name).
		Line().
		Type().
		Id(name).
		Struct()
}

func (g *generator) VisitOr(n *model.Or) {
	name := nameOf(n)
	var members = "<TODO>"

	g.
		Commentf("%s is a union of %s.", name, members).
		Line().
		Type().
		Id(name).
		Interface()
}

func (g *generator) VisitTuple(n *model.Tuple) {
	name := nameOf(n)

	g.
		Commentf("%s is a %d-tuple.", name, len(n.Types)).
		Line().
		Type().
		Id(name).
		Struct()
}

func (g *generator) VisitStructLit(n *model.StructLit) {
	name := nameOf(n)

	g.
		Commentf("%s is a literal structure.", name).
		Line().
		Type().
		Id(name).
		Struct()
}

func (g *generator) VisitStringLit(n *model.StringLit) {}
