package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

var tupleNames = map[int]string{
	2: "Pair",
	3: "Triad",
	4: "Tetrad",
	5: "Pentad",
}

func tupleName(arity int) string {
	name, ok := tupleNames[arity]
	if !ok {
		panic(fmt.Sprintf("no name for tuple with arity %d", arity))
	}

	return fmt.Sprintf("%sOf", name)
}

func (g *generator) tupleTypeExpr(t *metamodel.Type) jen.Code {
	array := true
	for _, item := range t.Items[1:] {
		if !reflect.DeepEqual(item, t.Items[0]) {
			array = false
			break
		}
	}

	length := len(t.Items)

	if array {
		return jen.
			Index(jen.Lit(length)).
			Add(g.typeExpr(t.Items[0]))
	}

	g.tupleArities[length] = struct{}{}

	return jen.
		Id(tupleName(length)).
		TypesFunc(func(gen *jen.Group) {
			for _, item := range t.Items {
				gen.Add(g.typeExpr(item))
			}
		})
}

func (g *generator) generateTuples(gen *jen.File) {
	if len(g.tupleArities) == 0 {
		return
	}

	generateBanner(gen, "TUPLES")

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

		gen.Commentf("%s is a tuple of %d values.", tupleName(arity), arity)
		gen.Type().
			Id(tupleName(arity)).
			Types(jen.List(types...).Any()).
			Struct(fields...)

		gen.Commentf("%s returns a tuple of %d values.", name, arity)
		gen.Func().
			Id(name).
			Types(jen.List(types...).Any()).
			Params(parameters...).
			Add(fullType).
			Block(
				jen.Return(
					jen.Add(fullType).
						Values(variables...),
				),
			)

		name += "P"
		gen.Commentf("%s returns a pointer to a tuple of %d values.", name, arity)
		gen.Func().
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

		gen.Line().
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
					jen.Id("marshal").
						Call(
							jen.
								Index(jen.Lit(arity)).
								Any().
								Values(fieldExpressions...),
						),
				),
			)

		gen.Line().
			Func().
			Params(jen.Id("x").Add(fullType)).
			Id("UnmarshalJSON").
			Params(
				jen.Id("data").Index().Byte(),
			).
			Params(
				jen.Error(),
			).
			BlockFunc(func(gen *jen.Group) {
				gen.Var().
					Id("elements").
					Index(jen.Lit(arity)).
					Qual("encoding/json", "RawMessage")

				gen.Line().
					If(
						jen.Err().
							Op(":=").
							Id("unmarshal").
							Call(
								jen.Id("data"),
								jen.Op("&").
									Id("elements"),
							),
						jen.Err().
							Op("!=").
							Nil(),
					).
					Block(
						jen.Return(
							jen.Err(),
						),
					)

				for i := 0; i < arity; i++ {
					gen.Line().
						If(
							jen.Err().
								Op(":=").
								Id("unmarshal").
								Call(
									jen.Id("elements").
										Index(jen.Lit(i)),
									jen.Op("&").
										Add(fieldExpressions[i]),
								),
							jen.Err().
								Op("!=").
								Nil(),
						).
						Block(
							jen.Return(
								jen.Err(),
							),
						)
				}

				gen.Line().
					Return(
						jen.Nil(),
					)
			})
	}
}
