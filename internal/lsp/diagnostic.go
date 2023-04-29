package lsp

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
)

func (h *handler) HandleTextDocumentDiagnostic(
	ctx context.Context,
	p proto.DocumentDiagnosticParams,
) (proto.DocumentDiagnosticReport, error) {
	rep := proto.DocumentDiagnosticReport{
		// First: &lsp.RelatedFullDocumentDiagnosticReport{
		// 	FullDocumentDiagnosticReport: lsp.FullDocumentDiagnosticReport{
		// 		Items: []lsp.Diagnostic{
		// 			{
		// 				Range: lsp.Range{
		// 					Start: lsp.Position{
		// 						Line:      10,
		// 						Character: 0,
		// 					},
		// 					End: lsp.Position{
		// 						Line:      10,
		// 						Character: 20,
		// 					},
		// 				},
		// 				Severity: lsp.DiagnosticSeverityError,
		// 				Message:  "It's all good, baby!",
		// 			},
		// 		},
		// 	},
		// },
	}

	return rep, nil
}
