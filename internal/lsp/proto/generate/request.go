package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) generateRequest(
	gen *jen.File,
	m metamodel.Request,
) {
	if m.Direction != "clientToServer" {
		return
	}

	name := normalizeName(m.Method)
	handlerName := name + "Handler"
	methodName := "Handle" + name
	routeName := name + "Route"

	gen.Commentf("%s handles %q requests.", handlerName, m.Method)
	gen.Type().
		Id(handlerName).
		InterfaceFunc(func(gen *jen.Group) {
			generateDocs(gen, m.Documentation)
			gen.Id(methodName).
				ParamsFunc(func(gen *jen.Group) {
					gen.Line().Qual("context", "Context")
					if m.Params != nil {
						gen.Line().Add(g.typeRef(m.Params))
					}
					gen.Line()
				}).
				ParamsFunc(func(gen *jen.Group) {
					if !m.Result.IsNull() {
						gen.Add(g.typeRef(m.Result))
					}
					gen.Error()
				})
		})

	gen.Commentf("%s returns a route for the %q request.", routeName, m.Method)
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
						handler := jen.Id("h").Dot(methodName)

						if m.Params == nil || m.Result.IsNull() {
							handler = jen.
								Func().
								ParamsFunc(func(gen *jen.Group) {
									gen.Id("ctx").Qual("context", "Context")

									if m.Params == nil {
										gen.Id("_").Struct()
									} else {
										gen.Id("p").Add(g.typeRef(m.Params))
									}
								}).
								ParamsFunc(func(gen *jen.Group) {
									if m.Result.IsNull() {
										gen.Any()
									} else {
										gen.Add(g.typeRef(m.Result))
									}

									gen.Error()
								}).
								Block(
									jen.ReturnFunc(func(gen *jen.Group) {
										if m.Result.IsNull() {
											gen.Nil()
										}

										gen.Add(handler).
											CallFunc(func(gen *jen.Group) {
												gen.Id("ctx")
												if m.Params != nil {
													gen.Id("p")
												}
											})
									}),
								)
						}

						gen.Line().Lit(m.Method)
						gen.Line().Add(handler)
						gen.Line()
					}),
			),
		)
}
