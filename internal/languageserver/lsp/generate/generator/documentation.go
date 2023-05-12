package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/wordwrap"
)

// documentation adds word-wrapped and sanitized documentation comments to a
// code block.
func documentation(
	code interface{ Comment(string) *jen.Statement },
	docs ...string,
) {
	text := strings.Join(docs, "\n\n")
	text = strings.TrimSpace(text)

	if text == "" {
		return
	}

	paragraphs := strings.Split(text, "\n\n")

	for pn, para := range paragraphs {
		if pn > 0 {
			code.Comment("")
		}

		para = strings.ReplaceAll(para, "\n-", "<list item>")
		para = strings.ReplaceAll(para, "\n", " ")
		para = strings.ReplaceAll(para, "<list item>", "\n-")

		lines := wordwrap.Wrap(para, 72)

		for _, line := range lines {
			if !strings.HasPrefix(line, "@since") {
				code.Comment(line)
			}
		}
	}
}
