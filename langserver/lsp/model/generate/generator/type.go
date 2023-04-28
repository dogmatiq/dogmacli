package generator

import (
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/metamodel"
)

func (g *generator) typeRef(t *metamodel.Type) jen.Code {
	switch t.Kind {
	case "base":
		return g.baseType(t)
	case "reference":
		return g.referenceType(t)
	case "tuple":
		return g.tupleType(t)
	case "or":
		return g.unionType(t)
	case "literal":
		return g.literalName(t)
	case "stringLiteral":
		// TODO:
		return jen.String()
	case "map":
		k := g.typeRef(t.MapKey)
		v := g.typeRef(t.MapValue())
		return jen.Map(k).Add(v)
	case "array":
		e := g.typeRef(t.ArrayElement)
		return jen.Index().Add(e)
	default:
		panic("unsupported kind: " + t.Kind)
	}
}

func (g *generator) baseType(t *metamodel.Type) jen.Code {
	switch t.Name {
	case "boolean":
		return jen.Bool()
	case "decimal":
		return jen.Float64()
	case "string":
		return jen.String()
	case "integer":
		return jen.Int32()
	case "uinteger":
		return jen.Uint32()
	case "DocumentUri":
		return jen.Id("DocumentURI")
	case "URI":
		return jen.Id("URI")
	case "null":
		return jen.Id("Null")
	default:
		panic("unsupported base type: " + t.Name)
	}
}

func (g *generator) referenceType(t *metamodel.Type) jen.Code {
	switch t.Name {
	case "LSPObject":
		k := jen.String()
		v := jen.Any()
		return jen.Map(k).Add(v)
	case "LSPArray":
		return jen.Index().Any()
	case "LSPAny":
		return jen.Any()
	default:
		return jen.Id(normalizeName(t.Name))
	}
}

func (g *generator) tupleType(t *metamodel.Type) jen.Code {
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

func (g *generator) unionType(t *metamodel.Type) jen.Code {
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
