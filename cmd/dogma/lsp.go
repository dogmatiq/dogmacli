package main

import (
	"os"

	"github.com/dogmatiq/dogmacli/internal/lsp"
	"github.com/dogmatiq/harpy"
	"github.com/dogmatiq/imbue"
	"go.uber.org/zap"
)

func init() {
	imbue.With0(
		catalog,
		func(
			ctx imbue.Context,
		) (*lsp.Server, error) {
			logger, err := zap.NewDevelopment(
				zap.WithCaller(false),
				zap.AddStacktrace(zap.FatalLevel), // disable stack trace for errors
			)
			if err != nil {
				return nil, err
			}

			return &lsp.Server{
				In:      os.Stdin,
				Out:     os.Stdout,
				Version: version,
				Logger: harpy.ZapExchangeLogger{
					Target: logger,
				},
			}, nil
		},
	)
}
