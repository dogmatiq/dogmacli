CGO_ENABLED=1

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/dogma
	$< $(RUN_ARGS)
