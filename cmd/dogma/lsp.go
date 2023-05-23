package main

import (
	"os"

	"github.com/dogmatiq/dogmacli/internal/langserver"
	"github.com/dogmatiq/harpy"
	"github.com/dogmatiq/imbue"
	"golang.org/x/exp/slog"
)

func init() {
	imbue.With0(
		catalog,
		func(
			ctx imbue.Context,
		) (*langserver.Server, error) {
			logger := slog.New(
				slog.NewTextHandler(
					os.Stderr,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)

			return &langserver.Server{
				In:             os.Stdin,
				Out:            os.Stdout,
				Version:        version,
				Logger:         logger,
				ExchangeLogger: harpy.NewSLogExchangeLogger(logger),
			}, nil
		},
	)
}
