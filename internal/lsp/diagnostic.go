package lsp

import (
	"context"
	"path/filepath"

	"github.com/dogmatiq/dogmacli/internal/linter"
	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
	"golang.org/x/tools/go/packages"
)

func (h *handler) HandleTextDocumentDiagnostic(
	ctx context.Context,
	p proto.DocumentDiagnosticParams,
) (proto.DocumentDiagnosticReport, error) {
	cfg := &packages.Config{
		Context: ctx,
		Dir:     filepath.Dir(p.TextDocument.URI.Path),
		Mode:    linter.PackageLoadMode,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return proto.DocumentDiagnosticReport{}, err
	}

	var rep *proto.RelatedFullDocumentDiagnosticReport

	for _, d := range linter.Lint(pkgs[0]) {
		if rep == nil {
			rep = &proto.RelatedFullDocumentDiagnosticReport{}
		}

		if d.Begin.Filename != p.TextDocument.URI.Path {
			continue
		}

		// TODO: the character position is _supposed_ to be measured in UTF-16
		// code units by default, but it appears to work correctly when
		// measuring bytes?
		rep.Items = append(
			rep.Items,
			proto.Diagnostic{
				Range: proto.Range{
					Start: proto.Position{
						Line:      uint32(d.Begin.Line - 1),
						Character: uint32(d.Begin.Column - 1),
					},
					End: proto.Position{
						Line:      uint32(d.End.Line - 1),
						Character: uint32(d.End.Column - 1),
					},
				},
				Severity: severity(d.Severity),
				Source:   "dogma",
				Message:  d.Message,
			},
		)
	}

	return proto.DocumentDiagnosticReport{
		First: rep,
	}, nil
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
