package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/metamodel"
)

// Generate returns generated code that represents the LSP data model.
func Generate(
	code *jen.File,
	root metamodel.Root,
) {
	v := declarationGenerator{File: code}

	for _, t := range root.NamedTypes {
		metamodel.VisitNamedType(t, v)
	}
}

type declarationGenerator struct {
	*jen.File
	typeExprGenerator
}

func (g declarationGenerator) VisitStructure(t metamodel.Structure) {
}

func (g declarationGenerator) VisitTypeAlias(t metamodel.TypeAlias) {
}
