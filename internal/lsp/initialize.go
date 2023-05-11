package lsp

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
)

func (h *handler) HandleInitialize(
	ctx context.Context,
	p proto.InitializeParams,
) (proto.InitializeResult, error) {
	for _, f := range p.WorkspaceFolders {
		h.addWorkspaceFolder(f)
	}

	return proto.InitializeResult{
		ServerInfo: &proto.InitializeResultServerInfo{
			Name:    proto.V("Dogma"),
			Version: proto.V(h.Version),
		},
		Capabilities: &proto.ServerCapabilities{
			TextDocumentSync: &proto.OneOf2[proto.TextDocumentSyncOptions, proto.TextDocumentSyncKind]{
				First: &proto.TextDocumentSyncOptions{
					OpenClose: proto.V(true),
					Change:    proto.TextDocumentSyncKindFull,
				},
			},
			Workspace: &proto.ServerCapabilitiesWorkspace{
				WorkspaceFolders: &proto.WorkspaceFoldersServerCapabilities{
					Supported: proto.V(false),
					ChangeNotifications: &proto.OneOf2[string, bool]{
						Second: proto.V(true),
					},
				},
			},
			DiagnosticProvider: &proto.OneOf2[proto.DiagnosticOptions, proto.DiagnosticRegistrationOptions]{
				First: &proto.DiagnosticOptions{
					Identifier:            proto.V("dogma"),
					InterFileDependencies: proto.V(true),
					WorkspaceDiagnostics:  proto.V(false),
				},
			},
		},
	}, nil
}
