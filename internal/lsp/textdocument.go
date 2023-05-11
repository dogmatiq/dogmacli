package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
)

func (h *handler) HandleTextDocumentDidOpen(
	ctx context.Context,
	p proto.DidOpenTextDocumentParams,
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
	p proto.DidChangeTextDocumentParams,
) error {
	before := string(h.overlay[p.TextDocument.URI.Path])
	after := before

	for _, change := range p.ContentChanges {
		switch {
		case change.First != nil:
			// TODO: UTF-16
			runes := []rune(after)
			start := positionToRuneOffset(runes, *change.First.Range.Start)
			end := positionToRuneOffset(runes, *change.First.Range.End)
			after = string(runes[:start]) + *change.First.Text + string(runes[end:])
		case change.Second != nil:
			after = *change.Second.Text
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
	p proto.DidCloseTextDocumentParams,
) error {
	// delete(h.overlay, p.TextDocument.URI.Path)
	return nil
}

func (h *handler) HandleTextDocumentDiagnostic(
	ctx context.Context,
	p proto.DocumentDiagnosticParams,
) (proto.DocumentDiagnosticReport, error) {
	var rep proto.RelatedFullDocumentDiagnosticReport

	for _, f := range h.workspaceFolders {
		for _, d := range f.Diagnostics[p.TextDocument.URI.Path] {
			// TODO: ast "column" is bytes, but LSP "character" is UTF-16 characters
			rep.Items = append(
				rep.Items,
				proto.Diagnostic{
					Range: &proto.Range{
						Start: &proto.Position{
							Line:      proto.V(uint32(d.Begin.Line - 1)),
							Character: proto.V(uint32(d.Begin.Column - 1)),
						},
						End: &proto.Position{
							Line:      proto.V(uint32(d.End.Line - 1)),
							Character: proto.V(uint32(d.End.Column - 1)),
						},
					},
					Severity: severity(d.Severity),
					Source:   proto.V("dogma"),
					Message:  proto.V(d.Message),
				},
			)
		}
	}

	x, err := json.Marshal(rep)
	fmt.Fprintln(os.Stderr, string(x), err)

	return proto.DocumentDiagnosticReport{First: &rep}, nil
}

func severity(s diagnostic.Severity) proto.DiagnosticSeverity {
	switch s {
	case diagnostic.Error:
		return proto.DiagnosticSeverityError
	case diagnostic.Warning:
		return proto.DiagnosticSeverityWarning
	case diagnostic.Improvement:
		return proto.DiagnosticSeverityInformation
	default:
		panic("unknown severity")
	}
}

func positionToRuneOffset(text []rune, p proto.Position) int {
	if *p.Line == 0 {
		return int(*p.Character)
	}

	lines := 0
	var prev rune

	for offset, char := range text {
		if lines == int(*p.Line) {
			return offset + int(*p.Character)
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
