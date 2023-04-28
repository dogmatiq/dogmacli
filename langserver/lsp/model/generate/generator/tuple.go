package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

func tupleName(arity int) string {
	name, ok := tupleNames[arity]
	if !ok {
		panic(fmt.Sprintf("no name for tuple with arity %d", arity))
	}

	return fmt.Sprintf("%sOf", name)
}

func (g *generator) tuples(code *jen.File) {
	for arity := 2; arity <= g.tupleArity; arity++ {
		var types []jen.Code
		var variables []jen.Code
		var fields []jen.Code

		for i := 0; i < arity; i++ {
			types = append(
				types,
				jen.Id(fmt.Sprintf("T%d", i)),
			)
			variables = append(
				variables,
				jen.Id(
					strings.ToLower(ordinalFieldNames[i]),
				),
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
				"%s is a tuple of %d values.",
				tupleName(arity),
				arity,
			)
		code.
			Type().
			Id(tupleName(arity)).
			Types(
				jen.List(types...).Any(),
			).
			StructFunc(func(code *jen.Group) {
				for i, f := range fields {
					code.
						Add(f).
						Add(types[i])
				}
			})

		funcName, ok := tupleNames[arity]
		if !ok {
			panic(fmt.Sprintf("no name for tuple with arity %d", arity))
		}

		code.
			Commentf(
				"%s returns a tuple of %d values.",
				funcName,
				arity,
			)
		code.
			Func().
			Id(funcName).
			Types(
				jen.List(types...).Any(),
			).
			ParamsFunc(func(code *jen.Group) {
				for i, v := range variables {
					code.
						Add(v).
						Add(types[i])
				}
			}).
			Id(tupleName(arity)).Types(types...).
			BlockFunc(func(code *jen.Group) {
				code.
					Return(
						jen.
							Id(tupleName(arity)).
							Types(types...).
							Values(variables...),
					)
			})

		funcName += "P"
		code.
			Commentf(
				"%s returns a pointer to a tuple of %d values.",
				funcName,
				arity,
			)
		code.
			Func().
			Id(funcName).
			Types(
				jen.List(types...).Any(),
			).
			ParamsFunc(func(code *jen.Group) {
				for i, v := range variables {
					code.
						Add(v).
						Add(types[i])
				}
			}).
			Op("*").Id(tupleName(arity)).Types(types...).
			BlockFunc(func(code *jen.Group) {
				code.
					Return(
						jen.
							Op("&").
							Id(tupleName(arity)).
							Types(types...).
							Values(variables...),
					)
			})
	}
}

var tupleNames = map[int]string{
	2: "Pair",
	3: "Triad",
	4: "Tetrad",
	5: "Pentad",
}
