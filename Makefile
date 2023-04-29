GO_EMBEDDED_FILES += internal/lsp/proto/metamodel/metamodel-3.17.0.json
DOGMACLI_GENERATED_FILES += internal/lsp/proto/lsp.gen.go

GENERATED_FILES += $(DOGMACLI_GENERATED_FILES)

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

run: $(GO_DEBUG_DIR)/dogma
	$< $(args)

$(DOGMACLI_GENERATED_FILES): $(shell find internal/lsp/proto/generate -type f)
	go run ./internal/lsp/proto/generate -- $@
