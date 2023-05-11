package lsp

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/linter"
	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

func (h *handler) lint(ctx context.Context, f *workspaceFolder) error {
	cfg := &packages.Config{
		Context: ctx,
		Dir:     f.Dir,
		Mode:    linter.PackageLoadMode,
		Overlay: h.overlay,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		h.Logger.Error(
			"unable to load packages for linting",
			slog.String("path", f.Dir),
			slog.String("error", err.Error()),
		)

		return nil
	}

	diagnostics := linter.Lint(pkgs...)
	f.Diagnostics = map[string][]diagnostic.Diagnostic{}

	for _, d := range diagnostics {
		h.Logger.Debug(
			"diagnostic",
			slog.String("file", d.Begin.Filename),
		)
		f.Diagnostics[d.Begin.Filename] = append(f.Diagnostics[d.Begin.Filename], d)
	}

	h.Logger.Debug(
		"linted package",
		slog.String("path", f.Dir),
		slog.Int("diagnostics", len(diagnostics)),
	)

	return nil
}
