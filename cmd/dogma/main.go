package main

import (
	"fmt"
	"os"

	"github.com/dogmatiq/dogmacli/cmd/dogma/internal/commands/visualize"
	"github.com/spf13/cobra"
)

// version string, automatically set during build process.
var version = "0.0.0"

var root = &cobra.Command{
	Use:     "dogma",
	Short:   "dogma command-line tools",
	Version: version,
}

func main() {
	root.AddCommand(visualize.Root)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
