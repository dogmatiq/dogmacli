package main

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/langserver/lsp/model/generate/generator"
	"golang.org/x/exp/slices"
)

func main() {
	i := slices.Index(os.Args, "--")
	args := os.Args[i+1:]
	filename := args[0]
	pkgPath := path.Join("github.com/dogmatiq/dogmacli/", path.Dir(filename))

	code := jen.NewFilePathName(pkgPath, path.Base(pkgPath))
	code.HeaderComment("Code generated by Dogma CLI's build process. DO NOT EDIT.")

	generator.Generate(code)

	buf := &bytes.Buffer{}
	if err := code.Render(buf); err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		filename,
		buf.Bytes(),
		0644,
	); err != nil {
		panic(fmt.Sprintf("unable to write to file: %s", err))
	}
}
