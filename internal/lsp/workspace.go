package lsp

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
	"golang.org/x/exp/slog"
)

type workspaceFolder struct {
	Dir         string
	Diagnostics map[string][]diagnostic.Diagnostic
}

func (h *handler) HandleWorkspaceDidChangeWorkspaceFolders(
	ctx context.Context,
	p proto.DidChangeWorkspaceFoldersParams,
) error {
	for _, f := range p.Event.Removed {
		h.removeWorkspaceFolder(f)
	}
	for _, f := range p.Event.Added {
		h.addWorkspaceFolder(f)
	}
	return nil
}

func (h *handler) addWorkspaceFolder(f proto.WorkspaceFolder) {
	if h.workspaceFolders == nil {
		h.workspaceFolders = map[proto.URI]*workspaceFolder{}
	}

	wf := &workspaceFolder{
		Dir: f.URI.Path,
	}

	h.workspaceFolders[*f.URI] = wf

	h.Logger.Debug(
		"added workspace folder",
		slog.String("name", *f.Name),
		slog.String("path", f.URI.Path),
	)
}

func (h *handler) removeWorkspaceFolder(f proto.WorkspaceFolder) {
	delete(h.workspaceFolders, *f.URI)

	h.Logger.Debug(
		"removed workspace folder",
		slog.String("name", *f.Name),
		slog.String("path", f.URI.Path),
	)
}
