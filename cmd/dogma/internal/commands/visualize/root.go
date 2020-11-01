package visualize

import "github.com/spf13/cobra"

// Root is the root command from this package.
var Root = &cobra.Command{
	Use:     "visualize",
	Short:   "generate visualizations of Dogma applications",
	Aliases: []string{"vis", "viz"},
}
