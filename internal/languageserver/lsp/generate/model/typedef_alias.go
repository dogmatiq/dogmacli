package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Alias describes a named type alias.
type Alias struct {
	TypeName      string
	Documentation Documentation
	Type          Type
}

// Name returns the type name.
func (d Alias) Name() string {
	return d.TypeName
}

func (b *builder) alias(in lowlevel.Alias, out *Alias) {
	out.TypeName = in.Name
	out.Documentation = in.Documentation
	out.Type = b.typeRef(in.Type)
}

func (d Alias) accept(v TypeDefVisitor) { v.Alias(d) }
func (v *typeDefX[T]) Alias(d Alias)    { v.V = v.X.Alias(d) }
