package commands

import (
	"os"

	"github.com/dogmatiq/dogmacli/cmd/dogma/internal/commands/lsp"
	"github.com/dogmatiq/dogmacli/cmd/dogma/internal/commands/visualize"
	"github.com/dogmatiq/imbue"
	"github.com/spf13/cobra"
)

// Root returns the root "dogma" command.
func Root(con *imbue.Container, ver string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "dogma",
		DisableFlagsInUseLine: true,
		Version:               ver,
		Short:                 "Dogma command-line",
	}

	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	cmd.AddCommand(
		lsp.Command(con),
		visualize.Command(con),
	)

	return cmd
}
