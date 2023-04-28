package generator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
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

func (g *generator) tupleRef(t *metamodel.Type) jen.Code {
	n := len(t.Items)

	first := t.Items[0]
	array := true
	for _, item := range t.Items[1:] {
		if !reflect.DeepEqual(item, first) {
			array = false
			break
		}
	}

	if array {
		return jen.Index(
			jen.Lit(n),
		).Add(
			g.typeRef(first),
		)
	}

	var types []jen.Code
	for _, item := range t.Items {
		types = append(types, g.typeRef(item))
	}

	if g.tupleArities == nil {
		g.tupleArities = map[int]struct{}{}
	}
	g.tupleArities[n] = struct{}{}

	return jen.
		Id(tupleName(n)).
		Types(types...)
}

func (g *generator) generateTuples(gen *jen.File) {
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

		gen.Line()
		gen.Func().
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
								Index(jen.Lit(arity)).
								Any().
								Values(fieldExpressions...),
						),
				),
			)

		gen.Line()
		gen.
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
					Qual(
						"encoding/json",
						"RawMessage",
					)

				gen.Line()
				gen.If(
					jen.Err().
						Op(":=").
						Qual(
							"encoding/json",
							"Unmarshal",
						).
						Call(
							jen.Id("data"),
							jen.Op("&").
								Id("elements"),
						),
					jen.Err().
						Op("!=").
						Nil(),
				).Block(
					jen.Return(
						jen.Err(),
					),
				)

				for i := 0; i < arity; i++ {
					gen.Line()
					gen.If(
						jen.Err().
							Op(":=").
							Qual(
								"encoding/json",
								"Unmarshal",
							).
							Call(
								jen.Id("elements").
									Index(jen.Lit(i)),
								jen.Op("&").
									Add(fieldExpressions[i]),
							),
						jen.Err().
							Op("!=").
							Nil(),
					).Block(
						jen.Return(
							jen.Err(),
						),
					)
				}

				gen.Line()
				gen.Return(
					jen.Nil(),
				)
			})
	}
}
