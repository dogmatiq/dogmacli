package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func unionName(arity int) string {
	return fmt.Sprintf("OneOf%d", arity)
}

func (g *generator) unionRef(t *metamodel.Type) jen.Code {
	if len(t.Items) == 2 {
		first, second := t.Items[0], t.Items[1]
		if first.Kind == "base" && first.Name == "null" {
			return jen.Op("*").Add(g.typeRef(second))
		}
		if second.Kind == "base" && second.Name == "null" {
			return jen.Op("*").Add(g.typeRef(first))
		}
	}

	var types []jen.Code
	for _, item := range t.Items {
		types = append(types, g.typeRef(item))
	}

	n := len(types)
	if g.unionArities == nil {
		g.unionArities = map[int]struct{}{}
	}
	g.unionArities[n] = struct{}{}

	return jen.
		Id(unionName(n)).
		Types(types...)
}

func (g *generator) generateUnions(gen *jen.File) {
	for _, arity := range sortedKeys(g.unionArities) {
		var (
			types            []jen.Code
			fields           []jen.Code
			fieldNames       []jen.Code
			fieldExpressions []jen.Code
		)

		for i := 0; i < arity; i++ {
			n := ordinalFieldNames[i]
			t := fmt.Sprintf("T%d", i)

			types = append(types, jen.Id(t))
			fields = append(fields, jen.Id(n).Op("*").Id(t))
			fieldNames = append(fieldNames, jen.Id(n))
			fieldExpressions = append(fieldExpressions, jen.Id("x").Dot(n))
		}

		fullType := jen.Id(unionName(arity)).Types(types...)

		gen.Commentf("%s is a union of %d values.", unionName(arity), arity)
		gen.Type().
			Id(unionName(arity)).
			Types(jen.List(types...).Any()).
			Struct(fields...)

		gen.Func().
			Params(jen.Id("x").Add(fullType)).
			Id("MarshalJSON").
			Params().
			Params(
				jen.Index().Byte(),
				jen.Error(),
			).
			BlockFunc(func(gen *jen.Group) {
				for _, f := range fieldExpressions {
					gen.If(jen.Add(f).Op("!=").Nil()).
						Block(
							jen.Return(
								jen.Qual("encoding/json", "Marshal").
									Call(f),
							),
						)
				}

				gen.Return(
					jen.Index().
						Byte().
						Call(
							jen.Lit("null"),
						),
					jen.Nil(),
				)
			})

		gen.Line()
		gen.Func().
			Params(jen.Id("x").Op("*").Add(fullType)).
			Id("UnmarshalJSON").
			Params(
				jen.Id("data").Index().Byte(),
			).
			Params(
				jen.Error(),
			).
			BlockFunc(func(gen *jen.Group) {
				gen.Op("*").
					Id("x").
					Op("=").
					Add(fullType).
					Values()

				gen.Line()
				gen.Var().
					Defs(
						jen.Err().Error(),
						jen.Id("errs").Index().Error(),
					)

				for _, f := range fieldExpressions {
					gen.Line()
					gen.Err().
						Op("=").
						Qual("encoding/json", "Unmarshal").
						Call(
							jen.Id("data"),
							jen.Op("&").
								Add(f),
						)

					gen.If(
						jen.Err().
							Op("==").
							Nil(),
					).Block(
						jen.Return(
							jen.Nil(),
						),
					)

					gen.Id("errs").
						Op("=").
						Append(
							jen.Id("errs"),
							jen.Err(),
						)

				}

				gen.Line()
				gen.Return(
					jen.Qual("errors", "Join").
						Call(
							jen.Id("errs").
								Op("..."),
						),
				)
			})
	}
}
