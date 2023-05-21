package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Alias describes a named type alias.
type Alias struct {
	typeDefNode

	UnderlyingType Type
}

func (b *builder) buildAlias(in lowlevel.Alias, out *Alias) {
	out.name = in.Name
	out.docs = in.Documentation

	out.UnderlyingType = b.buildType(in.Type)
}
