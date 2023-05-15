package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

// Generate returns generated code that represents the LSP data model.
func Generate(
	m model.Model,
	f *jen.File,
) {
	g := &generator{
		Model: m,
		File:  f,
	}

	g.Generate()
}

type generator struct {
	Model model.Model
	File  *jen.File
}

func (g *generator) Generate() {
	v := &typeDefGen{
		File: g.File,
	}
	for _, d := range g.Model.TypeDefs {
		model.VisitTypeDef(d, v)
		g.File.Line()
	}
}
