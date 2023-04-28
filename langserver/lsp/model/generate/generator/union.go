package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func unionName(arity int) string {
	return fmt.Sprintf("OneOf%d", arity)
}

func (g *generator) unions(code *jen.File) {
	for arity := 2; arity <= g.unionArity; arity++ {
		var types []jen.Code
		var fields []jen.Code

		for i := 0; i < arity; i++ {
			types = append(
				types,
				jen.Id(fmt.Sprintf("T%d", i)),
			)
			fields = append(
				fields,
				jen.Id(
					ordinalFieldNames[i],
				),
			)
		}

		code.
			Commentf(
				"%s is a union of %d values.",
				unionName(arity),
				arity,
			)
		code.
			Type().
			Id(unionName(arity)).
			Types(
				jen.List(types...).Any(),
			).
			StructFunc(func(code *jen.Group) {
				for i, f := range fields {
					code.
						Add(f).
						Op("*").
						Add(types[i])
				}
			})

		code.
			Func().
			Params(
				jen.
					Id("v").
					Id(unionName(arity)).
					Types(types...),
			).
			Id("MarshalJSON").
			Params().
			Params(
				jen.Index().Byte(),
				jen.Error(),
			).
			BlockFunc(func(code *jen.Group) {
				for _, f := range fields {
					code.
						If(
							jen.
								Id("v").
								Op(".").
								Add(f).
								Op("!=").
								Nil(),
						).
						BlockFunc(func(code *jen.Group) {
							code.
								Return(
									jen.
										Qual(
											"encoding/json",
											"Marshal",
										).
										Call(
											jen.
												Id("v").
												Op(".").
												Add(f),
										),
								)
						})
				}

				code.
					Return(
						jen.Index().Byte().Call(
							jen.Lit("null"),
						),
						jen.Nil(),
					)
			})

		code.Line()

		code.
			Func().
			Params(
				jen.
					Id("v").
					Op("*").
					Id(unionName(arity)).
					Types(types...),
			).
			Id("UnmarshalJSON").
			Params(
				jen.Id("data").Index().Byte(),
			).
			Params(
				jen.Error(),
			).
			BlockFunc(func(code *jen.Group) {
				code.
					Var().
					Defs(
						jen.Err().Error(),
						jen.Id("errs").Index().Error(),
					)

				code.Line()

				for _, f := range fields {
					code.
						Id("v").
						Op(".").
						Add(f).
						Op("=").
						Nil()
				}

				code.Line()

				for _, f := range fields {
					code.
						Err().
						Op("=").
						Qual(
							"encoding/json",
							"Unmarshal",
						).
						Call(
							jen.Id("data"),
							jen.
								Op("&").
								Id("v").
								Op(".").
								Add(f),
						)

					code.
						If(
							jen.
								Err().
								Op("==").
								Nil(),
						).
						Block(
							jen.
								Return(
									jen.Nil(),
								),
						)

					code.
						Id("errs").
						Op("=").
						Append(
							jen.Id("errs"),
							jen.Err(),
						)

					code.Line()
				}

				code.
					Return(
						jen.
							Qual(
								"errors",
								"Join",
							).
							Call(
								jen.
									Id("errs").
									Op("..."),
							),
					)
			})
	}
}
