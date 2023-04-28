package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/generate/metamodel"
)

func (g *generator) isOmittable(t *metamodel.Type) bool {
	switch t.Kind {
	case "base":
		return true
	case "reference":
		if strings.HasPrefix(t.Name, "LSP") {
			return true
		}
		return g.names[normalizeName(t.Name)].omittable
	case "tuple":
		return false
	case "or":
		return false
	case "literal":
		return false
	case "stringLiteral":
		return false
	case "map":
		return true
	case "array":
		return true
	default:
		panic("unsupported kind: " + t.Kind)
	}
}

func (g *generator) typeRef(t *metamodel.Type) jen.Code {
	switch t.Kind {
	case "base":
		return g.baseRef(t)
	case "reference":
		return g.namedRef(t)
	case "tuple":
		return g.tupleRef(t)
	case "or":
		return g.unionRef(t)
	case "literal":
		return g.literalStructRef(t)
	case "stringLiteral":
		return g.literalStringRef(t)
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

func (g *generator) baseRef(t *metamodel.Type) jen.Code {
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

func (g *generator) namedRef(t *metamodel.Type) jen.Code {
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
