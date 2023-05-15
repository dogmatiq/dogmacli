GO_EMBEDDED_FILES += internal/languageserver/lsp/metamodel/internal/lowlevel/metamodel-3.17.0.json
GENERATED_FILES += internal/languageserver/lsp/lsp.gen.go

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

run: $(GO_DEBUG_DIR)/dogma
	$< $(args)

internal/languageserver/lsp/lsp.gen.go: $(shell find internal/languageserver/lsp/generate -type f)
	go run ./internal/languageserver/lsp/generate -- $@
