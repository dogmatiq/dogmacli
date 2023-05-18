package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

// documentation adds word-wrapped and sanitized documentation comments to a
// code block.
func documentation(
	code interface{ Comment(string) *jen.Statement },
	docs model.Documentation,
	format string,
	args ...any,
) {
	text := docs.Text + "\n\n" + fmt.Sprintf(format, args...)
	text = strings.TrimSpace(text)

	for _, line := range strings.Split(text, "\n") {
		code.Comment(line)
	}

	if docs.DeprecationMessage != "" {
		code.Comment("")
		code.Comment("Deprecated: " + docs.DeprecationMessage)
	}
}
