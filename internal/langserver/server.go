package langserver

import (
	"context"
	"io"

	"github.com/dogmatiq/dogmacli/internal/langserver/lsp"
	"github.com/dogmatiq/dogmacli/internal/langserver/transport"
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
	// h := &handler{
	// 	Version: s.Version,
	// 	Logger:  s.Logger,
	// }

	r := harpy.NewRouter(
	// lsp.InitializeRoute(h),
	// lsp.TextDocumentDiagnosticRoute(h),
	// lsp.WorkspaceDidChangeWorkspaceFoldersRoute(h),
	// lsp.TextDocumentDidOpenRoute(h),
	// lsp.TextDocumentDidChangeRoute(h),
	// lsp.TextDocumentDidCloseRoute(h),
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
	workspaceFolders map[lsp.URI]*workspaceFolder
}
