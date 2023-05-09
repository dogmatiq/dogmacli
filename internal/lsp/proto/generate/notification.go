package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateNotifications(gen *jen.File) {
	generateBanner(gen, "NOTIFICATIONS")

	for _, m := range g.root.Notifications {
		if m.Direction == "clientToServer" {
			g.generateNotification(gen, m)
		}
	}
}

func (g *generator) generateNotification(
	gen *jen.File,
	m metamodel.Notification,
) {
	name := normalizeName(m.Method)
	handlerName := name + "Handler"
	methodName := "Handle" + name
	routeName := name + "Route"

	gen.Commentf("%s handles %q notifications.", handlerName, m.Method)
	gen.Type().
		Id(handlerName).
		InterfaceFunc(func(gen *jen.Group) {
			generateDocs(gen, m.Documentation)
			gen.Id(methodName).
				ParamsFunc(func(gen *jen.Group) {
					gen.Line().Qual("context", "Context")
					if m.Params != nil {
						gen.Line().Add(g.typeExpr(m.Params))
					}
					gen.Line()
				}).
				ParamsFunc(func(gen *jen.Group) {
					gen.Error()
				})
		})

	gen.Commentf("%s returns a route for the %q notification.", routeName, m.Method)
	gen.Func().
		Id(routeName).
		ParamsFunc(func(gen *jen.Group) {
			gen.Id("h").Id(handlerName)
		}).
		ParamsFunc(func(gen *jen.Group) {
			gen.Qual("github.com/dogmatiq/harpy", "RouterOption")
		}).
		Block(
			jen.Return(
				jen.Qual("github.com/dogmatiq/harpy", "WithRoute").
					CallFunc(func(gen *jen.Group) {
						handler := jen.
							Func().
							ParamsFunc(func(gen *jen.Group) {
								gen.Id("ctx").Qual("context", "Context")

								if m.Params == nil {
									gen.Id("_").Struct()
								} else {
									gen.Id("p").Add(g.typeExpr(m.Params))
								}
							}).
							Params(
								jen.Any(),
								jen.Error(),
							).
							Block(
								jen.ReturnFunc(func(gen *jen.Group) {
									gen.Nil()
									gen.Id("h").Dot(methodName).
										CallFunc(func(gen *jen.Group) {
											gen.Id("ctx")
											if m.Params != nil {
												gen.Id("p")
											}
										})
								}),
							)

						gen.Line().Lit(m.Method)
						gen.Line().Add(handler)
						gen.Line()
					}),
			),
		)
}
