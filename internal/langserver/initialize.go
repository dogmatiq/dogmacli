package langserver

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/langserver/lsp"
)

func (h *handler) HandleInitialize(
	ctx context.Context,
	p lsp.InitializeParams,
) (lsp.InitializeResult, error) {
	if folders, ok := p.WorkspaceFolders.Get(); ok {
		if folders, ok := folders.(lsp.WorkspaceFolderArray); ok {
			for _, f := range folders {
				if err := h.addWorkspaceFolder(ctx, f); err != nil {
					return lsp.InitializeResult{}, err
				}
			}
		}
	}

	return lsp.InitializeResult{
		ServerInfo: lsp.With(lsp.InitializeResultServerInfo{
			Name:    "Dogma",
			Version: lsp.With(lsp.String(h.Version)),
		}),
		Capabilities: lsp.ServerCapabilities{
			PositionEncoding: lsp.With(lsp.UTF16PositionEncodingKind),
			TextDocumentSync: lsp.With[lsp.ServerCapabilitiesTextDocumentSync](
				lsp.TextDocumentSyncOptions{
					OpenClose: lsp.With(lsp.True),
					Change:    lsp.With(lsp.FullTextDocumentSyncKind),
				},
			),
			Workspace: lsp.With(lsp.ServerCapabilitiesWorkspace{
				WorkspaceFolders: lsp.With(lsp.WorkspaceFoldersServerCapabilities{
					Supported:           lsp.With(lsp.False),
					ChangeNotifications: lsp.With[lsp.WorkspaceFoldersServerCapabilitiesChangeNotifications](lsp.True),
				}),
			}),
			DiagnosticProvider: lsp.With[lsp.ServerCapabilitiesDiagnosticProvider](
				lsp.DiagnosticOptions{
					Identifier:            lsp.With(lsp.String("dogma")),
					InterFileDependencies: lsp.True,
					WorkspaceDiagnostics:  lsp.False,
				},
			),
		},
	}, nil
}
