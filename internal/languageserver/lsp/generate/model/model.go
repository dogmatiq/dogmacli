package model

import (
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"
)

// Model is the root of the model.
type Model struct {
	Version  string
	TypeDefs []TypeDef
	Methods  []Method
}

// Documentation is a container for documentation-related meta-data.
type Documentation = lowlevel.Documentation

// Root returns the root node of the meta-model.
func Root() *Model {
	b := &builder{
		aliases: map[string]*Alias{},
		structs: map[string]*Struct{},
		enums:   map[string]*Enum{},
	}
	return b.BuildModel(lowlevel.Root())
}

type builder struct {
	aliases map[string]*Alias
	enums   map[string]*Enum
	structs map[string]*Struct
}

func (b *builder) BuildModel(in lowlevel.Model) *Model {
	out := &Model{
		Version: in.MetaData.Version,
	}

	// Pre-construct pointers to the all of the TypeDef instances so that they
	// may be referenced before they have been fully populated.
	for _, i := range in.Aliases {
		o := &Alias{}
		b.aliases[i.Name] = o
		out.TypeDefs = append(out.TypeDefs, o)
	}
	for _, i := range in.Enums {
		o := &Enum{}
		b.enums[i.Name] = o
		out.TypeDefs = append(out.TypeDefs, o)
	}
	for _, i := range in.Structs {
		o := &Struct{}
		b.structs[i.Name] = o
		out.TypeDefs = append(out.TypeDefs, o)
	}

	// Populate the TypeDef implementations.
	for _, d := range in.Aliases {
		b.alias(d, b.aliases[d.Name])
	}
	for _, d := range in.Enums {
		b.enum(d, b.enums[d.Name])
	}
	for _, d := range in.Structs {
		b.structure(d, b.structs[d.Name])
	}

	// Populate the methods.
	for _, m := range in.Requests {
		out.Methods = append(out.Methods, b.call(m))
	}
	for _, m := range in.Notifications {
		out.Methods = append(out.Methods, b.notification(m))
	}

	return out
}
