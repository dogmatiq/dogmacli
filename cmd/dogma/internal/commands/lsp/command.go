package lsp

import (
	"github.com/dogmatiq/imbue"
	"github.com/spf13/cobra"
)

// Command returns the "lsp" command.
func Command(con *imbue.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp",
		Short: "Listen for language server protocol requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
