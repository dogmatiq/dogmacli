package langserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dogmatiq/dogmacli/internal/langserver/lsp"
	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
)

func (h *handler) HandleTextDocumentDidOpen(
	ctx context.Context,
	p lsp.DidOpenTextDocumentParams,
) error {
	data, err := os.ReadFile(p.TextDocument.URI.Path)
	if err != nil {
		return nil
	}

	if h.overlay == nil {
		h.overlay = map[string][]byte{}
	}
	h.overlay[p.TextDocument.URI.Path] = data

	return nil
}

func (h *handler) HandleTextDocumentDidChange(
	ctx context.Context,
	p lsp.DidChangeTextDocumentParams,
) error {
	before := string(h.overlay[p.TextDocument.URI.Path])
	after := before

	for _, change := range p.ContentChanges {
		switch change := change.(type) {
		case lsp.TextDocumentContentChangeEventA:
			// TODO: UTF-16
			runes := []rune(after)
			start := positionToRuneOffset(runes, change.Range.Start)
			end := positionToRuneOffset(runes, change.Range.End)
			after = string(runes[:start]) + string(change.Text) + string(runes[end:])

		case lsp.TextDocumentContentChangeEventB:
			after = string(change.Text)
		}
	}

	h.overlay[p.TextDocument.URI.Path] = []byte(after)

	for _, f := range h.workspaceFolders {
		_, err := filepath.Rel(f.Dir, p.TextDocument.URI.Path)
		if err != nil {
			h.Logger.Debug(err.Error())
			continue
		}

		if err := h.lint(ctx, f); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) HandleTextDocumentDidClose(
	ctx context.Context,
	p lsp.DidCloseTextDocumentParams,
) error {
	// delete(h.overlay, p.TextDocument.URI.Path)
	return nil
}

func (h *handler) HandleTextDocumentDiagnostic(
	ctx context.Context,
	p lsp.DocumentDiagnosticParams,
) (lsp.DocumentDiagnosticReport, error) {
	var rep lsp.RelatedFullDocumentDiagnosticReport

	for _, f := range h.workspaceFolders {
		for _, d := range f.Diagnostics[p.TextDocument.URI.Path] {
			// TODO: ast "column" is bytes, but LSP "character" is UTF-16 characters
			rep.Items = append(
				rep.Items,
				lsp.Diagnostic{
					Range: lsp.Range{
						Start: lsp.Position{
							Line:      lsp.UInt(d.Begin.Line - 1),
							Character: lsp.UInt(d.Begin.Column - 1),
						},
						End: lsp.Position{
							Line:      lsp.UInt(d.End.Line - 1),
							Character: lsp.UInt(d.End.Column - 1),
						},
					},
					Severity: severity(d.Severity),
					Source:   "dogma",
					Message:  lsp.String(d.Message),
				},
			)
		}
	}

	x, err := json.Marshal(rep)
	fmt.Fprintln(os.Stderr, string(x), err)

	return rep, nil
}

func severity(s diagnostic.Severity) lsp.DiagnosticSeverity {
	switch s {
	case diagnostic.Error:
		return lsp.ErrorDiagnosticSeverity
	case diagnostic.Warning:
		return lsp.WarningDiagnosticSeverity
	case diagnostic.Improvement:
		return lsp.InformationDiagnosticSeverity
	default:
		panic("unknown severity")
	}
}

func positionToRuneOffset(text []rune, p lsp.Position) int {
	if p.Line == 0 {
		return int(p.Character)
	}

	lines := 0
	var prev rune

	for offset, char := range text {
		if lines == int(p.Line) {
			return offset + int(p.Character)
		}

		switch {
		case char == '\r':
			lines++
		case char == '\n' && prev != '\r':
			lines++
		}

		prev = char
	}

	return len(text)
}
