package generator

import (
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

type typeInfo struct {
	Name        string
	Kind        reflect.Kind
	UseOptional bool

	expr  *jen.Statement
	reify func(typeInfo)
}

func (i typeInfo) IsReified() bool {
	return i.reify != nil
}

func (i typeInfo) Expr() *jen.Statement {
	if i.reify != nil {
		i.reify(i)
	}
	if i.expr != nil {
		return i.expr
	}
	return jen.Id(i.Name)
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
	return typeInfo{
		Name: "Bool",
		Kind: reflect.Bool,
	}
}

func (g *typeInfoX) Decimal() typeInfo {
	return typeInfo{
		Name:        "Decimal",
		Kind:        reflect.Float64,
		UseOptional: true,
	}
}

func (g *typeInfoX) String() typeInfo {
	return typeInfo{
		Name: "String",
		Kind: reflect.String,
	}
}

func (g *typeInfoX) Integer() typeInfo {
	return typeInfo{
		Name:        "Integer",
		Kind:        reflect.Int32,
		UseOptional: true,
	}
}

func (g *typeInfoX) UInteger() typeInfo {
	return typeInfo{
		Name:        "UInteger",
		Kind:        reflect.Int32,
		UseOptional: true,
	}
}

func (g *typeInfoX) DocumentURI() typeInfo {
	return typeInfo{
		Name:        "DocumentURI",
		Kind:        reflect.Pointer,
		UseOptional: true,
		expr:        jen.Op("*").Id("DocumentURI"),
	}
}

func (g *typeInfoX) URI() typeInfo {
	return typeInfo{
		Name:        "URI",
		Kind:        reflect.Pointer,
		UseOptional: true,
		expr:        jen.Op("*").Id("URI"),
	}
}

func (g *typeInfoX) Null() typeInfo {
	return typeInfo{}
}

func (g *typeInfoX) Array(t model.Array) typeInfo {
	return typeInfo{
		Name: fmt.Sprintf(
			"%sArray",
			g.typeInfo(t.Element).Name,
		),
		Kind:        reflect.Slice,
		UseOptional: false,
		reify:       func(i typeInfo) { g.reifyType(i.Name, t) },
	}
}

func (g *typeInfoX) Map(t model.Map) typeInfo {
	return typeInfo{
		Name: fmt.Sprintf(
			"%s%sMap",
			g.typeInfo(t.Key).Name,
			g.typeInfo(t.Value).Name,
		),
		Kind:        reflect.Map,
		UseOptional: false,
		reify:       func(i typeInfo) { g.reifyType(i.Name, t) },
	}
}

func (g *typeInfoX) And(t model.And) typeInfo {
	return typeInfo{
		Name:        g.nameFromScope(),
		Kind:        reflect.Struct,
		UseOptional: true,
		reify:       func(i typeInfo) { g.reifyType(i.Name, t) },
	}
}

func (g *typeInfoX) Or(t model.Or) typeInfo {
	return typeInfo{
		Name:  g.nameFromScope(),
		Kind:  reflect.Interface,
		reify: func(i typeInfo) { g.reifyType(i.Name, t) },
	}
}

func (g *typeInfoX) Tuple(t model.Tuple) typeInfo {
	return typeInfo{
		Name:        g.nameFromScope(),
		Kind:        reflect.Struct,
		UseOptional: true,
		reify:       func(i typeInfo) { g.reifyType(i.Name, t) },
	}
}

func (g *typeInfoX) StructLit(t model.StructLit) typeInfo {
	return typeInfo{
		Name:        g.nameFromScope(),
		Kind:        reflect.Struct,
		UseOptional: true,
		reify:       func(i typeInfo) { g.reifyType(i.Name, t) },
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
		Name:        name,
		Kind:        underlying.Kind,
		UseOptional: underlying.UseOptional,
	}
}

func (g *typeInfoX) Enum(d model.Enum) typeInfo {
	name := identifier(d.TypeName)
	underlying := g.typeInfo(d.Type)

	return typeInfo{
		Name:        name,
		Kind:        underlying.Kind,
		UseOptional: true,
	}
}

func (g *typeInfoX) Struct(d model.Struct) typeInfo {
	name := identifier(d.TypeName)

	return typeInfo{
		Name:        name,
		Kind:        reflect.Struct,
		UseOptional: true,
	}
}
