package main

import (
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto/metamodel"
)

func (g *generator) isOmittable(t *metamodel.Type) bool {
	switch t.Kind {
	case "base":
		return true
	case "reference":
		if strings.HasPrefix(t.Name, "LSP") {
			return true
		}
		for _, m := range g.root.Enumerations {
			if m.Name == t.Name {
				return g.isOmittable(m.Type)
			}
		}
		for _, m := range g.root.TypeAliases {
			if m.Name == t.Name {
				return g.isOmittable(m.Type)
			}
		}
		return false
	case "tuple":
		return false
	case "or":
		t, _ := normalizeUnion(t)
		if t.Kind == "or" {
			return false
		}
		return g.isOmittable(t)
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

func (g *generator) typeExpr(t *metamodel.Type) jen.Code {
	switch t.Kind {
	case "base":
		return g.baseTypeExpr(t)
	case "reference":
		return g.refTypeExpr(t)
	case "tuple":
		return g.tupleTypeExpr(t)
	case "or":
		return g.unionTypeExpr(t)
	case "literal":
		return g.literalStructTypeExpr(t)
	case "stringLiteral":
		return g.literalStringTypeExpr(t)
	case "map":
		return g.mapTypeExpr(t)
	case "array":
		return g.arrayTypeExpr(t)
	default:
		panic("unsupported kind: " + t.Kind)
	}
}

func (g *generator) baseTypeExpr(t *metamodel.Type) jen.Code {
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
}

func (g *generator) refTypeExpr(t *metamodel.Type) jen.Code {
	switch t.Name {
	case "LSPObject":
		return jen.Map(jen.String()).Add(jen.Any())
	case "LSPArray":
		return jen.Index().Any()
	case "LSPAny":
		return jen.Any()
	default:
		return jen.Id(normalizeName(t.Name))
	}
}

func (g *generator) mapTypeExpr(t *metamodel.Type) jen.Code {
	return jen.Id("Map").Types(
		g.typeExpr(t.MapKey),
		g.typeExpr(t.MapValue()),
	)
}

func (g *generator) arrayTypeExpr(t *metamodel.Type) jen.Code {
	return jen.Id("Array").Types(
		g.typeExpr(t.ArrayElement),
	)
}

func (g *generator) hasValidateMethod(t *metamodel.Type) bool {
	switch t.Kind {
	case "base":
		return false
	case "reference":
		if strings.HasPrefix(t.Name, "LSP") {
			return false
		}
		for _, m := range g.root.TypeAliases {
			if m.Name == t.Name {
				return g.hasValidateMethod(m.Type)
			}
		}
	case "or":
		t, _ := normalizeUnion(t)
		if t.Kind != "or" {
			return g.hasValidateMethod(t)
		}
	}

	return true
}

func (g *generator) zeroValue(t *metamodel.Type) *jen.Statement {
	switch t.Kind {
	case "base":
		switch t.Name {
		case "boolean":
			return jen.False()
		case "string":
			return jen.Lit("")
		case "decimal", "integer", "uinteger":
			return jen.Lit(0)
		}
	case "reference":
		if strings.HasPrefix(t.Name, "LSP") {
			return jen.Nil()
		}
		for _, m := range g.root.Enumerations {
			if m.Name == t.Name {
				return g.zeroValue(m.Type)
			}
		}
		for _, m := range g.root.TypeAliases {
			if m.Name == t.Name {
				return g.zeroValue(m.Type)
			}
		}
	case "or":
		t, _ := normalizeUnion(t)
		if t.Kind != "or" {
			return g.zeroValue(t)
		}
	case "stringLiteral":
		return jen.Lit("")
	case "map", "array":
		return jen.Nil()
	}

	expr := g.typeExpr(t)
	return jen.Parens(
		jen.Add(expr).Values(),
	)
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
