package main

import (
	"fmt"
	"os"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/generator"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
	"golang.org/x/exp/slices"
)

func main() {
	i := slices.Index(os.Args, "--")
	args := os.Args[i+1:]
	filename := args[0]

	file := jen.NewFile("lsp")
	file.HeaderComment("Code generated by Dogma CLI. DO NOT EDIT.")

	generator.Generate(
		model.Get(),
		file,
	)

	if err := file.Save(filename); err != nil {
		fmt.Println("unable to save file:", err)
		os.Exit(1)
	}
}