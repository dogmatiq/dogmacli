package generator

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

type typeInfo struct {
	IsNamed            bool
	IsReified          bool
	IsNativelyOptional bool
	HasGoType          bool
}

func (g *Generator) typeInfo(t model.Type) typeInfo {
	return model.TypeTo[typeInfo](
		t,
		typeInfoX{},
	)
}

type typeInfoX struct{}

func (x typeInfoX) Bool() typeInfo {
	return typeInfo{
		IsNamed:            true,
		IsNativelyOptional: true,
		HasGoType:          true,
	}
}

func (x typeInfoX) Decimal() typeInfo {
	return typeInfo{
		IsNamed:   true,
		HasGoType: true,
	}
}

func (x typeInfoX) String() typeInfo {
	return typeInfo{
		IsNamed:            true,
		IsNativelyOptional: true,
		HasGoType:          true,
	}
}

func (x typeInfoX) Integer() typeInfo {
	return typeInfo{
		IsNamed:   true,
		HasGoType: true,
	}
}

func (x typeInfoX) UInteger() typeInfo {
	return typeInfo{
		IsNamed:   true,
		HasGoType: true,
	}
}

func (x typeInfoX) DocumentURI() typeInfo {
	return typeInfo{
		IsNamed:            true,
		IsNativelyOptional: true,
		HasGoType:          true,
	}
}

func (x typeInfoX) URI() typeInfo {
	return typeInfo{
		IsNamed:            true,
		IsNativelyOptional: true,
		HasGoType:          true,
	}
}

func (x typeInfoX) Null() typeInfo {
	return typeInfo{}
}

func (x typeInfoX) Array(t model.Array) typeInfo {
	return typeInfo{
		IsNativelyOptional: true,
		HasGoType:          true,
	}
}

func (x typeInfoX) Map(t model.Map) typeInfo {
	return typeInfo{
		IsNativelyOptional: true,
		HasGoType:          true,
	}
}

func (x typeInfoX) And(t model.And) typeInfo {
	return typeInfo{
		IsNamed:   true,
		IsReified: true,
		HasGoType: true,
	}
}

func (x typeInfoX) Or(t model.Or) typeInfo {
	return typeInfo{
		IsNamed:   true,
		IsReified: true,
		HasGoType: true,
	}
}

func (x typeInfoX) Tuple(t model.Tuple) typeInfo {
	return typeInfo{
		IsNamed:   true,
		IsReified: true,
		HasGoType: true,
	}
}

func (x typeInfoX) StructLit(t model.StructLit) typeInfo {
	return typeInfo{
		IsNamed:   true,
		IsReified: true,
		HasGoType: true,
	}
}

func (x typeInfoX) StringLit(t model.StringLit) typeInfo {
	return typeInfo{}
}

func (x typeInfoX) Reference(t model.Reference) typeInfo {
	return model.TypeDefTo[typeInfo](t.Target, x)
}

func (x typeInfoX) Alias(d model.Alias) typeInfo {
	info := model.TypeTo[typeInfo](d.Type, x)
	info.IsNamed = true
	return info
}

func (x typeInfoX) Enum(d model.Enum) typeInfo {
	return typeInfo{
		IsNamed:   true,
		HasGoType: true,
	}
}

func (x typeInfoX) Struct(d model.Struct) typeInfo {
	return typeInfo{
		IsNamed:   true,
		HasGoType: true,
	}
}
