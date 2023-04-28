package lsp

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/dogmacli/langserver"
	"github.com/dogmatiq/imbue"
	"github.com/spf13/cobra"
)

// Command returns the "lsp" command.
func Command(con *imbue.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp",
		Short: "Listen for language server protocol requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			stdio, err := cmd.Flags().GetBool("stdio")
			if err != nil {
				return err
			}

			if !stdio {
				return errors.New("--stdio is currently the only supported mode")
			}

			cmd.SilenceUsage = true

			ctx, cancel := signal.NotifyContext(
				cmd.Context(),
				os.Interrupt,
				syscall.SIGTERM,
			)
			defer cancel()

			return imbue.Invoke1(
				ctx,
				con,
				func(
					ctx context.Context,
					s *langserver.Server,
				) error {
					result := make(chan error, 1)
					go func() {
						result <- s.Run(ctx)
					}()

					select {
					case err := <-result:
						return err
					case <-ctx.Done():
						// Don't wait for s.Run() to finish when context is
						// canceled.
						return nil
					}
				},
			)
		},
	}

	cmd.Flags().Bool("stdio", false, "use stdio for communication")

	return cmd
}
