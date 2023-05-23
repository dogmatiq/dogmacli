GO_EMBEDDED_FILES += internal/langserver/lsp/metamodel/internal/lowlevel/metamodel-3.17.0.json
GENERATED_FILES += internal/langserver/lsp/lsp.gen.go

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

run: $(GO_DEBUG_DIR)/dogma
	$< $(args)

internal/langserver/lsp/lsp.gen.go: $(shell find internal/langserver/lsp/generate -type f)
	go run ./internal/langserver/lsp/generate -- $@
