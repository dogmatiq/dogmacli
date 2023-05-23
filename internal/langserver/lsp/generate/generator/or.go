package generator

import (
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
	"golang.org/x/exp/slices"
)

func (g *generator) VisitOr(n *model.Or) {
	name := nameOf(n)

	// Find any other "or" types that n is a member of, so that we can also
	// implement its interface.
	var methods []string
	for _, t := range g.Model.Types {
		if t, ok := t.(*model.Or); ok {
			for _, member := range t.Types {
				if memberName, ok := tryNameOf(member); ok {
					if memberName == name {
						methods = append(methods, "is"+nameOf(t))
					}
				}
			}
		}
	}

	slices.Sort(methods)

	// Always place n's method first.
	methods = append([]string{"is" + name}, methods...)

	g.Commentf("%s is a union (aka 'or') of several other types.", name)
	g.Comment("")
	g.Comment("It may be one of the following types:")

	for _, t := range n.Types {
		if name, ok := tryNameOf(t); ok {
			g.Commentf("  - %s", name)
		}
	}

	g.
		Type().
		Id(name).
		InterfaceFunc(
			g.withGroup(func() {
				for _, m := range methods {
					g.Id(m).Params()
				}
			}),
		)

	g.
		Func().
		Id("unmarshal"+name).
		Params(
			jen.Id("data").Index().Byte(),
			jen.Id("v").Id("*"+name),
		).
		Params(
			jen.Error(),
		).
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("err").Error()

			for _, t := range n.Types {
				if kindOf(t) == reflect.Invalid {
					continue
				}

				fn := jen.Qual("encoding/json", "Unmarshal")
				if kindOf(t) == reflect.Interface {
					fn = jen.Id("unmarshal" + nameOf(t))
				}

				g.Block(
					jen.
						Var().
						Id("candidate").
						Id(nameOf(t)),
					jen.Err().
						Op(":=").
						Add(fn).
						Call(
							jen.Id("data"),
							jen.Op("&").Id("candidate"),
						),
					jen.
						If(
							jen.Err().Op("==").Nil(),
						).
						Block(
							jen.Op("*").Id("v").Op("=").Id("candidate"),
							jen.Return(jen.Nil()),
						),
					jen.
						Err().Op("=").Qual("errors", "Join").
						Call(
							jen.Id("err"),
							jen.Err(),
						),
				)
			}

			g.Return(jen.Id("err"))
		})

	// Add methods to any (non-interface) types that are members of this
	// "or" type, such that they satisfy this interface.
	var members []string

	for _, member := range n.Types {
		switch kindOf(member) {
		case reflect.Interface:
		case reflect.Invalid:
		default:
			members = append(members, nameOf(member))
		}
	}

	slices.Sort(members)

	for _, member := range members {
		g.Line()

		for _, method := range methods {
			g.
				Func().
				Params(
					jen.Id(member),
				).
				Id(method).
				Params().
				Block()
		}
	}
}
