package visualize

import (
	"context"
	"io"
	"os"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/visualization/dot"
	"github.com/dogmatiq/pluginkit"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "graph <plugin> [...]",
		Short: "generate a visualization of one or more Dogma applications in Graphviz DOT format",
		Args:  cobra.MinimumNArgs(1),
		RunE:  graph,
	}

	cmd.Flags().StringP("output", "o", "-", "write output to the specified file")

	Root.AddCommand(cmd)
}

// graph is the entry point for the "graph" command.
func graph(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	applications, err := loadConfigsFromPlugins(ctx, args)
	if err != nil {
		return err
	}

	g, err := dot.Generate(applications...)
	if err != nil {
		return err
	}

	out, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	if out != "-" {
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		defer f.Close()

		w = f
	}

	_, err = io.WriteString(
		w,
		g.String(),
	)

	return err
}

// loadConfigsFromPlugins returns the configuration for all applications
// provided by the given plugin files.
func loadConfigsFromPlugins(
	ctx context.Context,
	files []string,
) ([]configkit.Application, error) {
	var applications []configkit.Application

	for _, f := range files {
		p, err := pluginkit.Load(ctx, f)
		if err != nil {
			return nil, err
		}
		defer p.Close()

		s, ok := p.ApplicationService()
		if !ok {
			continue
		}

		for _, k := range s.ApplicationKeys() {
			app, closer, err := s.NewApplication(ctx, k)
			if err != nil {
				return nil, err
			}
			defer closer.Close()

			applications = append(
				applications,
				configkit.FromApplication(app),
			)
		}
	}

	return applications, nil
}
