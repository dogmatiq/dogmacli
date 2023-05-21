package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Model is the root of the model.
type Model struct {
	node

	Version string
	Defs    map[string]Def
	Types   []Type
}

// Get returns the root node of the meta-model.
func Get() *Model {
	b := &builder{}
	return b.build(lowlevel.Get())
}

func (b *builder) build(in lowlevel.Model) *Model {
	b.model = &Model{
		Version: in.MetaData.Version,
		Defs:    map[string]Def{},
	}
	b.parent = b.model

	// Pre-construct pointers to the all of the Def instances so that they may
	// be referenced before they have been fully populated.
	for _, def := range in.Aliases {
		buildDef(b, def.Name, def, b.buildAlias)
	}
	for _, def := range in.Enums {
		buildDef(b, def.Name, def, b.buildEnum)
	}
	for _, def := range in.Structs {
		buildDef(b, def.Name, def, b.buildStruct)
	}
	for _, def := range in.Requests {
		buildDef(b, def.Method, def, b.buildCall)
	}
	for _, def := range in.Notifications {
		buildDef(b, def.Method, def, b.buildNotification)
	}

	for _, fn := range b.resolvers {
		fn()
	}

	return b.model
}
