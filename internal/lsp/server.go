package lsp

import (
	"context"
	"io"

	"github.com/dogmatiq/dogmacli/internal/lsp/transport"
	"github.com/dogmatiq/harpy"
)

// Server is a JSON-RPC server that handles LSP requests.
type Server struct {
	In      io.ReadCloser
	Out     io.Writer
	Version string
	Logger  harpy.ExchangeLogger
}

// Run starts the JSON-RPC server.
func (s *Server) Run(ctx context.Context) error {
	return transport.Run(
		ctx,
		newExchanger(s.Version),
		s.In,
		s.Out,
		s.Logger,
	)
}
