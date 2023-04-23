package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/dogmacli/cmd/dogma/internal/commands"
	"github.com/dogmatiq/imbue"
)

// version string, automatically set during build process.
var version = "0.0.0"

// catalog is the dependency injection catalog for the Grit CLI.
var catalog = imbue.NewCatalog()

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	con := imbue.New(imbue.WithCatalog(catalog))
	defer con.Close()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	return commands.
		Root(con, version).
		ExecuteContext(ctx)
}
