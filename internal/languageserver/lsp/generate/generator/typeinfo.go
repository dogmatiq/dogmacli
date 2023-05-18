package generator

import (
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

type typeInfo struct {
	Name        *string
	NameHint    string
	TypeExpr    func() *jen.Statement
	TypeKind    reflect.Kind
	UseOptional bool
	IsReified   bool
}

func (g *Generator) typeInfo(t model.Type) typeInfo {
	return model.TypeTo[typeInfo](t, &typeInfoX{g})
}

func (g *Generator) typeInfoForDef(d model.TypeDef) typeInfo {
	return model.TypeDefTo[typeInfo](d, &typeInfoX{g})
}

type typeInfoX struct {
	*Generator
}

func (g *typeInfoX) Bool() typeInfo {
	name := "Bool"
	return typeInfo{
		Name:     &name,
		NameHint: name,
		TypeExpr: func() *jen.Statement { return jen.Id(name) },
		TypeKind: reflect.Bool,
	}
}

func (g *typeInfoX) Decimal() typeInfo {
	name := "Decimal"
	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    reflect.Float64,
		UseOptional: true,
	}
}

func (g *typeInfoX) String() typeInfo {
	name := "String"
	return typeInfo{
		Name:     &name,
		NameHint: name,
		TypeExpr: func() *jen.Statement { return jen.Id(name) },
		TypeKind: reflect.String,
	}
}

func (g *typeInfoX) Integer() typeInfo {
	name := "Integer"
	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    reflect.Int32,
		UseOptional: true,
	}
}

func (g *typeInfoX) UInteger() typeInfo {
	name := "UInteger"
	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    reflect.Int32,
		UseOptional: true,
	}
}

func (g *typeInfoX) DocumentURI() typeInfo {
	name := "DocumentURI"
	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    reflect.Pointer,
		UseOptional: true,
	}
}

func (g *typeInfoX) URI() typeInfo {
	name := "URI"
	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    reflect.Pointer,
		UseOptional: true,
	}
}

func (g *typeInfoX) Null() typeInfo {
	return typeInfo{}
}

func (g *typeInfoX) Array(t model.Array) typeInfo {
	e := g.typeInfo(t.Element)

	return typeInfo{
		NameHint:    fmt.Sprintf("%sArray", e.NameHint),
		TypeExpr:    func() *jen.Statement { return jen.Index().Add(e.TypeExpr()) },
		TypeKind:    reflect.Slice,
		UseOptional: false,
	}
}

func (g *typeInfoX) Map(t model.Map) typeInfo {
	k := g.typeInfo(t.Key)
	v := g.typeInfo(t.Value)

	return typeInfo{
		NameHint:    fmt.Sprintf("%s%sMap", k.NameHint, v.NameHint),
		TypeExpr:    func() *jen.Statement { return jen.Map(k.TypeExpr()).Add(v.TypeExpr()) },
		TypeKind:    reflect.Map,
		UseOptional: false,
	}
}

func (g *typeInfoX) And(t model.And) typeInfo {
	name := g.nameFromScope()

	return typeInfo{
		Name:     &name,
		NameHint: name,
		TypeExpr: func() *jen.Statement {
			g.reifyType(name, t)
			return jen.Id(name)
		},
		TypeKind:    reflect.Struct,
		UseOptional: true,
		IsReified:   true,
	}
}

func (g *typeInfoX) Or(t model.Or) typeInfo {
	name := g.nameFromScope()

	return typeInfo{
		Name:     &name,
		NameHint: name,
		TypeExpr: func() *jen.Statement {
			g.reifyType(name, t)
			return jen.Id(name)
		},
		TypeKind:    reflect.Struct,
		UseOptional: true,
		IsReified:   true,
	}
}

func (g *typeInfoX) Tuple(t model.Tuple) typeInfo {
	name := g.nameFromScope()

	return typeInfo{
		Name:     &name,
		NameHint: name,
		TypeExpr: func() *jen.Statement {
			g.reifyType(name, t)
			return jen.Id(name)
		},
		TypeKind:    reflect.Struct,
		UseOptional: true,
		IsReified:   true,
	}
}

func (g *typeInfoX) StructLit(t model.StructLit) typeInfo {
	name := g.nameFromScope()

	return typeInfo{
		Name:     &name,
		NameHint: name,
		TypeExpr: func() *jen.Statement {
			g.reifyType(name, t)
			return jen.Id(name)
		},
		TypeKind:    reflect.Struct,
		UseOptional: true,
		IsReified:   true,
	}
}

func (g *typeInfoX) StringLit(t model.StringLit) typeInfo {
	return typeInfo{}
}

func (g *typeInfoX) Reference(t model.Reference) typeInfo {
	return model.TypeDefTo[typeInfo](t.Target, g)
}

func (g *typeInfoX) Alias(d model.Alias) typeInfo {
	name := identifier(d.TypeName)
	underlying := g.typeInfo(d.Type)

	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    underlying.TypeKind,
		UseOptional: underlying.UseOptional,
	}
}

func (g *typeInfoX) Enum(d model.Enum) typeInfo {
	name := identifier(d.TypeName)
	underlying := g.typeInfo(d.Type)

	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    underlying.TypeKind,
		UseOptional: true,
	}
}

func (g *typeInfoX) Struct(d model.Struct) typeInfo {
	name := identifier(d.TypeName)

	return typeInfo{
		Name:        &name,
		NameHint:    name,
		TypeExpr:    func() *jen.Statement { return jen.Id(name) },
		TypeKind:    reflect.Struct,
		UseOptional: true,
	}
}
