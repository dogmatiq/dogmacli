package model

import (
	"math"

	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Enum describes a named enumeration type.
type Enum struct {
	typeDef

	Type    Type
	Lax     bool
	Members []*EnumMember
}

// EnumMember describes a value within an enumeration.
type EnumMember struct {
	node

	Name          string
	Documentation Documentation
	Value         any
}

// AcceptVisitor dispatches to the appropriate method on the given visitor.
func (n *EnumMember) AcceptVisitor(v Visitor) {
	v.VisitEnumMember(n)
}

func (b *builder) buildEnum(in lowlevel.Enum, out *Enum) {
	out.Documentation = in.Documentation
	out.Type = b.buildType(in.Type)
	out.Lax = in.SupportsCustomValues

	for _, m := range in.Members {
		out.Members = append(
			out.Members,
			build(b, func(n *EnumMember) {
				n.Name = m.Name
				n.Documentation = m.Documentation
				n.Value = normalizeEnumValue(m.Value)
			}),
		)
	}
}

func normalizeEnumValue(v any) any {
	if f, ok := v.(float64); ok {
		if math.Mod(f, 1) == 0 {
			return int(f)
		}
	}
	return v
}
