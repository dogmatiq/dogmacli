package visualize

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/static"
	"github.com/dogmatiq/configkit/visualization/dot"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

func init() {
	cmd := &cobra.Command{
		Use:   "graph [<package> ...]",
		Short: "generate a visualization of one or more Dogma applications in Graphviz DOT format",
		RunE:  graph,
	}

	cmd.Flags().StringP(
		"output", "o",
		"-",
		"write output to the specified file",
	)

	Root.AddCommand(cmd)
}

// graph is the entry point for the "graph" command.
func graph(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if len(args) == 0 {
		args = []string{"./..."}
	}

	apps, err := loadConfigsFromPackages(ctx, args)
	if err != nil {
		return err
	}

	if len(apps) == 0 {
		return errors.New("no applications found")
	}

	g, err := dot.Generate(apps...)
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

// loadConfigsFromPackages returns the configuration for all applications
// defined within the packages that match the given patterns.
func loadConfigsFromPackages(
	ctx context.Context,
	patterns []string,
) ([]configkit.Application, error) {
	var applications []configkit.Application

	for _, pattern := range patterns {
		cfg := packages.Config{
			Context: ctx,
			Mode:    packages.LoadAllSyntax,
			Dir:     pattern,
		}

		if filepath.Base(pattern) == "..." {
			cfg.Dir = filepath.Dir(pattern)
			pattern = "./..."
		} else {
			pattern = "."
		}

		pkgs, err := packages.Load(&cfg, pattern)
		if err != nil {
			return nil, err
		}

		for _, pkg := range pkgs {
			for _, err := range pkg.Errors {
				return nil, err
			}
		}

		applications = append(
			applications,
			static.FromPackages(pkgs)...,
		)

	}

	return applications, nil
}
