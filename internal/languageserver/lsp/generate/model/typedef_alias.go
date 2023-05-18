package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Alias describes a named type alias.
type Alias struct {
	typeDef

	Type Type
}

func (b *builder) buildAlias(in lowlevel.Alias, out *Alias) {
	out.Documentation = in.Documentation
	out.Type = b.buildType(in.Type)
}
