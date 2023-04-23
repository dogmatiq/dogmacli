package visualize

import (
	"github.com/dogmatiq/dogmacli/cmd/dogma/internal/commands/visualize/graph"
	"github.com/dogmatiq/imbue"
	"github.com/spf13/cobra"
)

// Command returns the "visualize" command.
func Command(con *imbue.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "visualize",
		Short:   "Generate visualizations of Dogma applications",
		Aliases: []string{"vis", "viz"},
	}

	cmd.AddCommand(
		graph.Command(con),
	)

	return cmd
}
