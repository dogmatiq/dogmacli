package generator

import (
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *Generator) emitMethod(m model.MethodDef) {
	// g.pushScope(m.MethodName())
	// model.VisitMethod(m, &method{g})
	// g.popScope()
}

type method struct{ *Generator }

func (g *method) Call(m *model.Call) {
	params := jen.Null()
	if m.ParamsType != nil {
		g.pushNestedScope("Params")
		info := g.typeInfo(m.ParamsType)
		if info.Kind != reflect.Invalid {
			params = info.Expr()
		}
		g.popNestedScope()
	}

	if m.RegistrationOptionsType != nil {
		g.pushNestedScope("RegistrationOptions")
		g.typeInfo(m.RegistrationOptionsType)
		g.popNestedScope()
	}

	result := jen.Null()
	if m.ResultType != nil {
		g.pushNestedScope("Result")
		info := g.typeInfo(m.ResultType)
		if info.Kind != reflect.Invalid {
			result = info.Expr()
		}
		g.popNestedScope()
	}

	if m.PartialResultType != nil {
		g.pushNestedScope("PartialResult")
		g.typeInfo(m.PartialResultType)
		g.popNestedScope()
	}

	if m.ErrorDataType != nil {
		g.pushNestedScope("ErrorData")
		g.typeInfo(m.ErrorDataType)
		g.popNestedScope()
	}

	if m.Direction() == model.HandledByLanguageServer {
		g.File.
			Type().
			Id(normalize(m.Name(), "Handler")).
			Interface(
				jen.
					Id(normalize("Handle", m.Name())).
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
			Id(normalize(m.Name())).
			ParamsFunc(func(grp *jen.Group) {
				grp.Id("ctx").Qual("context", "Context")
				if m.ParamsType != nil {
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

func (g *method) Notification(m *model.Notification) {
	params := jen.Null()
	if m.ParamsType != nil {
		g.pushNestedScope("Params")
		info := g.typeInfo(m.ParamsType)
		if info.Kind != reflect.Invalid {
			params = info.Expr()
		}
		g.popNestedScope()
	}

	if m.RegistrationOptionsType != nil {
		g.pushNestedScope("RegistrationOptions")
		g.typeInfo(m.RegistrationOptionsType)
		g.popNestedScope()
	}

	if m.Direction() == model.HandledByLanguageServer {
		g.File.
			Type().
			Id(normalize(m.Name(), "Handler")).
			Interface(
				jen.
					Id(normalize("Handle", m.Name())).
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
			Id(normalize(m.Name())).
			ParamsFunc(func(grp *jen.Group) {
				grp.Id("ctx").Qual("context", "Context")
				if m.ParamsType != nil {
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
