package langserver

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/langserver/lsp"
	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"golang.org/x/exp/slog"
)

type workspaceFolder struct {
	Dir         string
	Diagnostics map[string][]diagnostic.Diagnostic
}

func (h *handler) HandleWorkspaceDidChangeWorkspaceFolders(
	ctx context.Context,
	p lsp.DidChangeWorkspaceFoldersParams,
) error {
	for _, f := range p.Event.Removed {
		if err := h.removeWorkspaceFolder(ctx, f); err != nil {
			return err
		}
	}

	for _, f := range p.Event.Added {
		if err := h.addWorkspaceFolder(ctx, f); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) addWorkspaceFolder(ctx context.Context, f lsp.WorkspaceFolder) error {
	if h.workspaceFolders == nil {
		h.workspaceFolders = map[lsp.URI]*workspaceFolder{}
	}

	wf := &workspaceFolder{
		Dir: f.URI.Path,
	}

	h.workspaceFolders[f.URI] = wf

	h.Logger.Debug(
		"added workspace folder",
		slog.String("name", string(f.Name)),
		slog.String("path", f.URI.Path),
	)

	return h.lint(ctx, wf)
}

func (h *handler) removeWorkspaceFolder(ctx context.Context, f lsp.WorkspaceFolder) error {
	delete(h.workspaceFolders, f.URI)

	h.Logger.Debug(
		"removed workspace folder",
		slog.String("name", string(f.Name)),
		slog.String("path", f.URI.Path),
	)

	return nil
}
