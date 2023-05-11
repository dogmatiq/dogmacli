package main

import (
	"os"

	"github.com/dogmatiq/dogmacli/internal/lsp"
	"github.com/dogmatiq/harpy"
	"github.com/dogmatiq/imbue"
	"golang.org/x/exp/slog"
)

func init() {
	imbue.With0(
		catalog,
		func(
			ctx imbue.Context,
		) (*lsp.Server, error) {
			logger := slog.New(
				slog.HandlerOptions{
					Level: slog.LevelDebug,
				}.NewTextHandler(os.Stderr),
			)

			return &lsp.Server{
				In:             os.Stdin,
				Out:            os.Stdout,
				Version:        version,
				Logger:         logger,
				ExchangeLogger: harpy.NewSLogExchangeLogger(logger),
			}, nil
		},
	)
}
