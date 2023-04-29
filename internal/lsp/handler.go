package lsp

import (
	"context"

	"github.com/dogmatiq/dogmacli/internal/lsp/proto"
	"github.com/dogmatiq/harpy"
)

func newExchanger(version string) harpy.Exchanger {
	h := &handler{
		version: version,
	}

	return harpy.NewRouter(
		proto.InitializeRoute(h),
		proto.TextDocumentDiagnosticRoute(h),
	)
}

type handler struct {
	version string
}

func (h *handler) HandleInitialize(
	ctx context.Context,
	p proto.InitializeParams,
) (proto.InitializeResult, error) {
	return proto.InitializeResult{
		ServerInfo: &proto.InitializeResultServerInfo{
			Name:    "Dogma",
			Version: h.version,
		},
		Capabilities: proto.ServerCapabilities{
			TextDocumentSync: &proto.OneOf2[proto.TextDocumentSyncOptions, proto.TextDocumentSyncKind]{
				First: &proto.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    proto.TextDocumentSyncKindFull,
				},
			},
			DiagnosticProvider: &proto.OneOf2[proto.DiagnosticOptions, proto.DiagnosticRegistrationOptions]{
				First: &proto.DiagnosticOptions{
					Identifier:            "dogma",
					InterFileDependencies: true,
					WorkspaceDiagnostics:  false,
				},
			},
		},
	}, nil
}
