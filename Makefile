GO_EMBEDDED_FILES += langserver/lsp/generate/metamodel/metamodel-3.17.0.json
DOGMACLI_GENERATED_FILES += langserver/lsp/model.gen.go

GENERATED_FILES += $(DOGMACLI_GENERATED_FILES)

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/dogma
	$< $(args)

$(DOGMACLI_GENERATED_FILES): $(shell find langserver/lsp/generate -type f)
	go run langserver/lsp/generate/main.go -- $@
