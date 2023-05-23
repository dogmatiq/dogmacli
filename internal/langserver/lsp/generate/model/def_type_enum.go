package model

import (
	"math"

	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model/internal/lowlevel"
)

// Enum describes a named enumeration type.
type Enum struct {
	typeDefNode

	UnderlyingType Type
	Members        []*EnumMember
	Strict         bool
}

// EnumMember describes a value within an enumeration.
type EnumMember struct {
	node

	Name          string
	Documentation Documentation
	Value         any
}

func (b *builder) buildEnum(in lowlevel.Enum, out *Enum) {
	out.name = in.Name
	out.docs = in.Documentation

	out.UnderlyingType = b.buildType(in.Type)

	for _, m := range in.Members {
		out.Members = append(out.Members, b.buildEnumMember(m))
	}

	out.Strict = !in.SupportsCustomValues
}

func (b *builder) buildEnumMember(in lowlevel.EnumMember) *EnumMember {
	return build(b, func(out *EnumMember) {
		out.Name = in.Name
		out.Documentation = in.Documentation
		out.Value = normalizeEnumValue(in.Value)
	})
}

func normalizeEnumValue(v any) any {
	if f, ok := v.(float64); ok {
		if math.Mod(f, 1) == 0 {
			return int(f)
		}
	}
	return v
}
