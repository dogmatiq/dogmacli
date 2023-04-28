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
	for _, arity := range sortedKeys(g.tupleArities) {
		name, ok := tupleNames[arity]
		if !ok {
			panic(fmt.Sprintf("no name for tuple with arity %d", arity))
		}

		var (
			types            []jen.Code
			variables        []jen.Code
			parameters       []jen.Code
			fields           []jen.Code
			fieldNames       []jen.Code
			fieldExpressions []jen.Code
		)

		for i := 0; i < arity; i++ {
			n := ordinalFieldNames[i]
			v := strings.ToLower(n)
			t := fmt.Sprintf("T%d", i)

			types = append(types, jen.Id(t))
			variables = append(variables, jen.Id(v))
			parameters = append(parameters, jen.Id(v).Id(t))
			fields = append(fields, jen.Id(n).Id(t))
			fieldNames = append(fieldNames, jen.Id(n))
			fieldExpressions = append(fieldExpressions, jen.Id("x").Dot(n))
		}

		fullType := jen.Id(tupleName(arity)).Types(types...)

		code.Commentf("%s is a tuple of %d values.", tupleName(arity), arity)
		code.
			Type().
			Id(tupleName(arity)).
			Types(jen.List(types...).Any()).
			Struct(fields...)

		code.Commentf("%s returns a tuple of %d values.", name, arity)
		code.
			Func().
			Id(name).
			Types(jen.List(types...).Any()).
			Params(parameters...).
			Add(fullType).
			Block(
				jen.Return(
					jen.
						Add(fullType).
						Values(variables...),
				),
			)

		name += "P"
		code.Commentf("%s returns a pointer to a tuple of %d values.", name, arity)
		code.
			Func().
			Id(name).
			Types(jen.List(types...).Any()).
			Params(parameters...).
			Op("*").Add(fullType).
			Block(
				jen.Return(
					jen.
						Op("&").
						Add(fullType).
						Values(variables...),
				),
			)

		code.Line()
		code.
			Func().
			Params(jen.Id("x").Add(fullType)).
			Id("MarshalJSON").
			Params().
			Params(
				jen.Index().Byte(),
				jen.Error(),
			).
			Block(
				jen.Return(
					jen.
						Qual(
							"encoding/json",
							"Marshal",
						).
						Call(
							jen.
								Index().
								Any().
								Values(fieldExpressions...),
						),
				),
			)
	}
}

var tupleNames = map[int]string{
	2: "Pair",
	3: "Triad",
	4: "Tetrad",
	5: "Pentad",
}
