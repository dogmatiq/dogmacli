package lsp

import (
	"context"
	"encoding/json"
	"os"

	"github.com/dogmatiq/dapper"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
)

func (s *Server) textDocumentDiagnostic(
	ctx context.Context,
	params proto.DocumentDiagnosticParams,
) (proto.DocumentDiagnosticReport, error) {
	dapper.Write(os.Stderr, params)

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

	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	enc.Encode(rep)

	return rep, nil
}
