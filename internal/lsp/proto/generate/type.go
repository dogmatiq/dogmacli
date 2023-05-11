package main

import (
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) typeExpr(t *metamodel.Type) jen.Code {
	switch t.Kind {
	case "base":
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
			panic("unexpected reference to null type")
		default:
			panic("unsupported base type: " + t.Name)
		}
	case "reference":
		switch t.Name {
		case "LSPObject":
			return jen.Id("Map").Types(
				jen.String(),
				jen.Any(),
			)
		case "LSPArray":
			return jen.Id("Array").Types(
				jen.Any(),
			)
		case "LSPAny":
			return jen.Any()
		default:
			return jen.Id(normalizeName(t.Name))
		}
	case "tuple":
		return g.tupleTypeExpr(t)
	case "or":
		return g.unionTypeExpr(t)
	case "literal":
		return g.literalStructTypeExpr(t)
	case "stringLiteral":
		return g.literalStringTypeExpr(t)
	case "map":
		return jen.Id("Map").Types(
			g.typeExpr(t.MapKey),
			g.typeExpr(t.MapValue()),
		)
	case "array":
		return jen.Id("Array").Types(
			g.typeExpr(t.ArrayElement),
		)
	default:
		panic("unsupported kind: " + t.Kind)
	}
}

type typeInfo struct {
	IsNillable     bool
	AddPointer     bool
	IsValidateable bool
}

func (g *generator) typeInfo(t *metamodel.Type) (info typeInfo) {
	switch t.Kind {
	case "base":
		return typeInfo{
			IsNillable:     false,
			AddPointer:     true,
			IsValidateable: false,
		}
	case "reference":
		switch t.Name {
		case "LSPObject":
			return typeInfo{
				IsNillable:     true,
				AddPointer:     false,
				IsValidateable: true,
			}
		case "LSPArray":
			return typeInfo{
				IsNillable:     true,
				AddPointer:     false,
				IsValidateable: true,
			}
		case "LSPAny":
			return typeInfo{
				IsNillable:     true,
				AddPointer:     false,
				IsValidateable: false,
			}
		}

		for _, m := range g.root.TypeAliases {
			if m.Name == t.Name {
				return g.typeInfo(m.Type)
			}
		}
		for _, m := range g.root.Enumerations {
			if m.Name == t.Name {
				return typeInfo{
					IsNillable:     false,
					AddPointer:     false,
					IsValidateable: true,
				}
			}
		}

		return typeInfo{
			IsNillable:     false,
			AddPointer:     true,
			IsValidateable: true,
		}
	case "tuple":
		return typeInfo{
			IsNillable:     false,
			AddPointer:     true,
			IsValidateable: true,
		}
	case "or":
		t, _ := normalizeUnion(t)
		if t.Kind != "or" {
			return g.typeInfo(t)
		}
		return typeInfo{
			IsNillable:     false,
			AddPointer:     true,
			IsValidateable: true,
		}
	case "literal":
		return typeInfo{
			IsNillable:     false,
			AddPointer:     true,
			IsValidateable: true,
		}
	case "stringLiteral":
		return typeInfo{
			IsNillable:     false,
			AddPointer:     false,
			IsValidateable: false,
		}
	case "map":
		return typeInfo{
			IsNillable:     true,
			AddPointer:     false,
			IsValidateable: true,
		}
	case "array":
		return typeInfo{
			IsNillable:     true,
			AddPointer:     false,
			IsValidateable: true,
		}
	default:
		panic("unsupported kind: " + t.Kind)
	}
}

func (g *generator) generateUniqueName(desired string) (actual, suffix string) {
	if _, ok := g.names[desired]; !ok {
		g.names[desired] = struct{}{}
		return desired, ""
	}

	i := 0
	for {
		i++

		suffix := strconv.Itoa(i)
		candidate := desired + suffix

		if _, ok := g.names[candidate]; !ok {
			g.names[candidate] = struct{}{}
			return candidate, suffix
		}
	}
}
