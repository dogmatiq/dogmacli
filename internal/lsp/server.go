package lsp

import (
	"context"
	"io"

	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
	"github.com/dogmatiq/dogmacli/internal/lsp/transport"
	"github.com/dogmatiq/harpy"
	"golang.org/x/exp/slog"
)

// Server is a JSON-RPC server that handles LSP requests.
type Server struct {
	In             io.ReadCloser
	Out            io.Writer
	Version        string
	Logger         *slog.Logger
	ExchangeLogger harpy.ExchangeLogger
}

// Run starts the JSON-RPC server.
func (s *Server) Run(ctx context.Context) error {
	h := &handler{
		Version: s.Version,
		Logger:  s.Logger,
	}

	r := harpy.NewRouter(
		proto.InitializeRoute(h),
		proto.TextDocumentDiagnosticRoute(h),
		proto.WorkspaceDidChangeWorkspaceFoldersRoute(h),
		proto.TextDocumentDidOpenRoute(h),
		proto.TextDocumentDidChangeRoute(h),
		proto.TextDocumentDidCloseRoute(h),
	)

	return transport.Run(
		ctx,
		r,
		s.In,
		s.Out,
		s.ExchangeLogger,
	)
}

// handler implements various generated LSP handler interfaces.
type handler struct {
	Version string
	Logger  *slog.Logger

	overlay          map[string][]byte
	workspaceFolders map[proto.URI]*workspaceFolder
}
