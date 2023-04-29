package lsp

import (
	"context"
	"io"

	lsp "github.com/dogmatiq/dogmacli/internal/lsp/proto"
	"github.com/dogmatiq/dogmacli/internal/lsp/transport"
	"github.com/dogmatiq/harpy"
)

type Server struct {
	In      io.ReadCloser
	Out     io.Writer
	Version string
	Logger  harpy.ExchangeLogger
}

func (s *Server) Run(ctx context.Context) error {
	return transport.Run(
		ctx,
		harpy.NewRouter(
			harpy.WithRoute(
				"initialize",
				s.initialize,
			),
			harpy.WithRoute(
				"textDocument/diagnostic",
				s.textDocumentDiagnostic,
			),
		),
		s.In,
		s.Out,
		s.Logger,
	)
}

func (s *Server) initialize(
	ctx context.Context,
	params lsp.InitializeParams,
) (lsp.InitializeResult, error) {
	return lsp.InitializeResult{
		ServerInfo: &lsp.InitializeResultServerInfo{
			Name:    "Dogma",
			Version: s.Version,
		},
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.OneOf2[lsp.TextDocumentSyncOptions, lsp.TextDocumentSyncKind]{
				First: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.TextDocumentSyncKindFull,
				},
			},
			DiagnosticProvider: &lsp.OneOf2[lsp.DiagnosticOptions, lsp.DiagnosticRegistrationOptions]{
				First: &lsp.DiagnosticOptions{
					Identifier:            "dogma",
					InterFileDependencies: true,
					WorkspaceDiagnostics:  false,
				},
			},
		},
	}, nil
}
