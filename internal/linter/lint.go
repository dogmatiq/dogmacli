package linter

import (
	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"github.com/dogmatiq/dogmacli/internal/linter/internal/lint"
	"github.com/dogmatiq/dogmacli/internal/linter/internal/rules/identity"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Lint runs the linter against the given packages.
func Lint(pkgs ...*packages.Package) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	prog, packages := ssautil.AllPackages(pkgs, SSABuilderMode)
	prog.Build()

	for i, pkg := range packages {
		if pkg == nil {
			continue
		}

		ctx := lint.NewContext(
			pkgs[i],
			prog,
			pkg,
		)

		identity.Lint(ctx)

		for _, d := range ctx.Diagnostics {
			diags = append(diags, *d)
		}
	}

	return diags
}

const (
	// PackageLoadMode is the packages.LoadMode required by the linter.
	PackageLoadMode = packages.NeedName |
		packages.NeedCompiledGoFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedTypesInfo |
		packages.NeedImports

	// SSABuilderMode is the ssa.BuilderMode required by the linter.
	SSABuilderMode = ssa.SanityCheckFunctions | ssa.GlobalDebug
)
