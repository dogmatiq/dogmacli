package metamodel

import (
	"math"

	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/metamodel/internal/lowlevel"
)

// Enum describes a named enumeration type.
type Enum struct {
	TypeName             string
	Documentation        Documentation
	Type                 Type
	Members              []EnumMember
	SupportsCustomValues bool
}

// EnumMember describes a value within an enumeration.
type EnumMember struct {
	Name          string
	Documentation Documentation
	Value         any
}

// Name returns the type name.
func (d Enum) Name() string {
	return d.TypeName
}

func (b *builder) enum(in lowlevel.Enum, out *Enum) {
	out.TypeName = in.Name
	out.Documentation = in.Documentation
	out.Type = b.typeRef(in.Type)
	out.SupportsCustomValues = in.SupportsCustomValues

	for _, m := range in.Members {
		out.Members = append(
			out.Members,
			EnumMember{
				Name:          m.Name,
				Documentation: m.Documentation,
				Value:         normalizeEnumValue(m.Value),
			},
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

func (d Enum) accept(v TypeDefVisitor) { v.Enum(d) }
func (v *typeDefX[T]) Enum(d Enum)     { v.V = v.X.Enum(d) }
