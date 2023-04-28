package langserver

import (
	"context"
	"io"

	"github.com/dogmatiq/dogmacli/langserver/iotransport"
	"github.com/dogmatiq/harpy"
)

type Server struct {
	In      io.ReadCloser
	Out     io.Writer
	Version string
	Logger  harpy.ExchangeLogger
}

func (s *Server) Run(ctx context.Context) error {
	return iotransport.Run(
		ctx,
		harpy.NewRouter(
		// harpy.WithRoute(
		// 	"initialize",
		// 	s.initialize,
		// 	harpy.AllowUnknownFields(true),
		// ),
		// harpy.WithRoute(
		// 	"textDocument/diagnostic",
		// 	s.textDocumentDiagnostic,
		// 	harpy.AllowUnknownFields(true),
		// ),
		),
		s.In,
		s.Out,
		s.Logger,
	)
}

// func (s *Server) initialize(
// 	ctx context.Context,
// 	params lsp.InitializeParams,
// ) (lsp.InitializeResult, error) {
// 	dapper.Write(os.Stderr, params)
// 	return lsp.InitializeResult{
// 		ServerInfo: lsp.ServerInfo{
// 			Name:    "Dogma",
// 			Version: s.Version,
// 		},
// 		Capabilities: lsp.ServerCapabilities{
// 			TextDocumentSync: lsp.TextDocumentSyncOptions{
// 				OpenClose: true,
// 				Change:    lsp.TextDocumentSyncKindFull,
// 			},
// 			DiagnosticProvider: lsp.DiagnosticRegistrationOptions{
// 				TextDocumentRegistrationOptions: lsp.TextDocumentRegistrationOptions{
// 					DocumentSelector: []lsp.DocumentFilter{},
// 				},
// 				DiagnosticOptions: lsp.DiagnosticOptions{
// 					Identifier:            "dogma",
// 					InterFileDependencies: true,
// 					WorkspaceDiagnostics:  false,
// 				},
// 				StaticRegistrationOptions: lsp.StaticRegistrationOptions{
// 					ID: "",
// 				},
// 			},
// 		},
// 	}, nil
// }
