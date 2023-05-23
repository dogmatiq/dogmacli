package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
)

// documentation adds word-wrapped and sanitized documentation comments to a
// code block.
func documentation(
	code interface{ Comment(string) *jen.Statement },
	docs model.Documentation,
	format string,
	args ...any,
) {
	text := fmt.Sprintf(format, args...)

	if s := strings.TrimSpace(docs.Text); s != "" {
		text += "\n\n"
		text += s
	}

	if s := strings.TrimSpace(docs.DeprecationMessage); s != "" {
		text += "\n\n"
		text += "Deprecated: " + s
	}

	text = strings.TrimSpace(text)
	for _, line := range strings.Split(text, "\n") {
		code.Comment("  " + line)
	}
}
