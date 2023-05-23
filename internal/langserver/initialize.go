package langserver

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/langserver/lsp"
)

func (h *handler) HandleInitialize(
	ctx context.Context,
	p lsp.InitializeParams,
) (lsp.InitializeResult, error) {
	for _, f := range p.WorkspaceFolders.(lsp.WorkspaceFolderArray) {
		if err := h.addWorkspaceFolder(ctx, f); err != nil {
			return lsp.InitializeResult{}, err
		}
	}

	return lsp.InitializeResult{
		ServerInfo: &lsp.InitializeResultServerInfo{
			Name:    "Dogma",
			Version: lsp.String(h.Version),
		},
		Capabilities: lsp.ServerCapabilities{
			PositionEncoding: lsp.UTF16PositionEncodingKind,
			TextDocumentSync: lsp.TextDocumentSyncOptions{
				OpenClose: lsp.True,
				Change:    lsp.FullTextDocumentSyncKind,
			},
			Workspace: &lsp.ServerCapabilitiesWorkspace{
				WorkspaceFolders: &lsp.WorkspaceFoldersServerCapabilities{
					Supported:           lsp.False,
					ChangeNotifications: lsp.True,
				},
			},
			DiagnosticProvider: lsp.DiagnosticOptions{
				Identifier:            "dogma",
				InterFileDependencies: lsp.True,
				WorkspaceDiagnostics:  lsp.False,
			},
		},
	}, nil
}
