package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) typeRef(t *metamodel.Type) jen.Code {
	switch t.Kind {
	case "reference":
		return jen.Id(normalizeName(t.Name))

	case "tuple":
		var types []jen.Code
		for _, item := range t.Items {
			types = append(types, g.typeRef(item))
		}

		n := len(types)
		if n > g.tupleArity {
			g.tupleArity = n
		}

		return jen.
			Id(tupleName(n)).
			Types(types...)

	case "or":
		var types []jen.Code
		for _, item := range t.Items {
			types = append(types, g.typeRef(item))
		}

		n := len(types)
		if n > g.unionArity {
			g.unionArity = n
		}

		return jen.
			Id(unionName(n)).
			Types(types...)

	case "literal":
		return g.literalName(t)

	case "stringLiteral":
		// TODO:
		return jen.String()

	case "map":
		return jen.Map(
			g.typeRef(t.MapKey),
		).Add(
			g.typeRef(t.MapValue()),
		)

	case "array":
		return jen.Index().Add(
			g.typeRef(t.ArrayElement),
		)

	case "base":
		return jen.Id(baseTypeName(t.Name))

	default:
		panic("unsupported kind: " + t.Kind)
	}
}

func baseTypeName(n string) string {
	switch n {
	case "boolean":
		return "bool"
	case "decimal":
		return "float64"
	case "string":
		return "string"
	case "integer":
		return "int32"
	case "uinteger":
		return "uint32"
	case "DocumentUri":
		return "DocumentURI"
	case "URI":
		return "URI"
	case "null":
		return "Null"
	default:
		panic("unsupported base type: " + n)
	}
}
