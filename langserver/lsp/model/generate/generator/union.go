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
						Add(types[i])
				}
			})
	}
}
