package langserver

import (
	"context"
	"encoding/json"
	"os"

	"github.com/dogmatiq/dapper"
	"github.com/dogmatiq/dogmacli/langserver/lsp"
)

func (s *Server) textDocumentDiagnostic(
	ctx context.Context,
	params lsp.DocumentDiagnosticParams,
) (lsp.DocumentDiagnosticReport, error) {
	dapper.Write(os.Stderr, params)

	rep := lsp.DocumentDiagnosticReport{
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

	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	enc.Encode(rep)

	return rep, nil
}
