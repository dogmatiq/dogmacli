package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *Generator) emitMethod(m model.Method) {
	g.pushScope(m.Name())
	model.VisitMethod(m, &method{g})
	g.popScope()
}

type method struct{ *Generator }

func (g *method) Call(m model.Call) {
	params := jen.Null()
	if m.Params != nil {
		g.pushNestedScope("Params")
		params = g.typeExpr(m.Params)
		g.popNestedScope()
	}

	if m.RegistrationOptions != nil {
		g.pushNestedScope("RegistrationOptions")
		g.typeExpr(m.RegistrationOptions)
		g.popNestedScope()
	}

	result := jen.Null()
	if m.Result != nil {
		g.pushNestedScope("Result")
		result = g.typeExpr(m.Result)
		g.popNestedScope()
	}

	if m.PartialResult != nil {
		g.pushNestedScope("PartialResult")
		g.typeExpr(m.PartialResult)
		g.popNestedScope()
	}

	if m.ErrorData != nil {
		g.pushNestedScope("ErrorData")
		g.typeExpr(m.ErrorData)
		g.popNestedScope()
	}

	if m.Direction == model.HandledByLanguageServer {
		g.File.
			Type().
			Id(identifier(m.MethodName, "Handler")).
			Interface(
				jen.
					Id(identifier("Handle", m.MethodName)).
					Params(
						jen.Qual("context", "Context"),
						params,
					).
					Params(
						result,
						jen.Error(),
					),
			)
	} else {
		g.File.
			Func().
			Params(
				jen.Id("c").Op("*").Id("Client"),
			).
			Id(identifier(m.MethodName)).
			ParamsFunc(func(grp *jen.Group) {
				grp.Id("ctx").Qual("context", "Context")
				if m.Params != nil {
					grp.Id("p").Add(params)
				}
			}).
			Params(
				result,
				jen.Error(),
			).
			Block(
				jen.Panic(
					jen.Lit("not implemented"),
				),
			)
	}
}

func (g *method) Notification(m model.Notification) {
	params := jen.Null()
	if m.Params != nil {
		g.pushNestedScope("Params")
		params = g.typeExpr(m.Params)
		g.popNestedScope()
	}

	if m.RegistrationOptions != nil {
		g.pushNestedScope("RegistrationOptions")
		g.typeExpr(m.RegistrationOptions)
		g.popNestedScope()
	}

	if m.Direction == model.HandledByLanguageServer {
		g.File.
			Type().
			Id(identifier(m.MethodName, "Handler")).
			Interface(
				jen.
					Id(identifier("Handle", m.MethodName)).
					Params(
						jen.Qual("context", "Context"),
						params,
					).
					Params(
						jen.Error(),
					),
			)
	} else {
		g.File.
			Func().
			Params(
				jen.Id("c").Op("*").Id("Client"),
			).
			Id(identifier(m.MethodName)).
			ParamsFunc(func(grp *jen.Group) {
				grp.Id("ctx").Qual("context", "Context")
				if m.Params != nil {
					grp.Id("p").Add(params)
				}
			}).
			Params(
				jen.Error(),
			).
			Block(
				jen.Panic(
					jen.Lit("not implemented"),
				),
			)
	}
}
