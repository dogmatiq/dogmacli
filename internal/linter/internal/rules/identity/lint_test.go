package identity_test

import (
	"testing"

	. "github.com/dogmatiq/dogmacli/internal/linter/internal/rules/identity"
	"github.com/dogmatiq/dogmacli/internal/linter/internal/rules/internal/ruletest"
)

func TestLint(t *testing.T) {
	ruletest.Run(t, Lint)
}
