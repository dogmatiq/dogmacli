package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Model is the root of the model.
type Model struct {
	node

	Version  string
	Methods  []Method
	TypeDefs []TypeDef
}

// Root returns the root node of the meta-model.
func Root() *Model {
	b := &builder{
		aliases: map[string]*Alias{},
		structs: map[string]*Struct{},
		enums:   map[string]*Enum{},
	}

	return b.buildModel(lowlevel.Root())
}

func (b *builder) buildModel(in lowlevel.Model) *Model {
	out := &Model{
		Version: in.MetaData.Version,
	}

	// Pre-construct pointers to the all of the TypeDef instances so that they
	// may be referenced before they have been fully populated.
	for _, a := range in.Aliases {
		n := &Alias{}
		n.name = a.Name
		n.setParent(out)

		b.aliases[a.Name] = n
		out.TypeDefs = append(out.TypeDefs, n)
	}
	for _, e := range in.Enums {
		n := &Enum{}
		n.name = e.Name
		n.setParent(out)

		b.enums[e.Name] = n
		out.TypeDefs = append(out.TypeDefs, n)
	}
	for _, s := range in.Structs {
		n := &Struct{}
		n.name = s.Name
		n.setParent(out)

		b.structs[s.Name] = n
		out.TypeDefs = append(out.TypeDefs, n)
	}

	// Populate the methods.
	for _, m := range in.Requests {
		out.Methods = append(out.Methods, b.buildCall(m))
	}
	for _, m := range in.Notifications {
		out.Methods = append(out.Methods, b.buildNotification(m))
	}

	// Populate the TypeDef implementations.
	for _, d := range in.Aliases {
		b.buildAlias(d, b.aliases[d.Name])
	}
	for _, d := range in.Enums {
		b.buildEnum(d, b.enums[d.Name])
	}
	for _, d := range in.Structs {
		b.buildStruct(d, b.structs[d.Name])
	}

	return out
}
