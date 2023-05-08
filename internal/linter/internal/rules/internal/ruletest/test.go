package ruletest

import (
	"context"
	"testing"

	linter "github.com/dogmatiq/dogmacli/internal/linter" // TODO: alias is to stop goimports from removing the import
	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"github.com/dogmatiq/dogmacli/internal/linter/internal/lint"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Rule is a linter rule.
type Rule func(*lint.Context)

// Run runs a linter rule against the "testdata" directory.
func Run(t *testing.T, rule Rule) {
	t.Helper()

	cfg := &packages.Config{
		Context: context.Background(),
		Dir:     "testdata",
		Mode:    linter.PackageLoadMode,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		t.Fatal(err)
	}

	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			t.Error(err)
		}
	}

	if t.Failed() {
		return
	}

	prog, spkgs := ssautil.AllPackages(pkgs, linter.SSABuilderMode)
	prog.Build()

	for i, pkg := range pkgs {
		ctx := lint.NewContext(
			pkg,
			prog,
			spkgs[i],
		)

		testPackage(t, ctx, rule)
	}
}

func testPackage(
	t *testing.T,
	ctx *lint.Context,
	rule Rule,
) {
	t.Helper()

	t.Run(ctx.Package.PkgPath, func(t *testing.T) {
		t.Helper()
		t.Parallel()

		expectations, errors := parseExpectations(ctx.Package)
		for _, err := range errors {
			t.Error(err)
		}

		if t.Failed() {
			return
		}

		rule(ctx)

		type failure struct {
			Expectation expectation
			Diagnostic  *diagnostic.Diagnostic
		}

		var (
			failures   []failure
			missing    = slices.Clone(expectations)
			unexpected = slices.Clone(ctx.Diagnostics)
		)

		for _, e := range expectations {
			for _, d := range ctx.Diagnostics {
				if e.isLocatedWithin(d) {
					remove(&missing, e)
					remove(&unexpected, d)

					if !e.isSatisfiedBy(d) {
						failures = append(failures, failure{e, d})
					}
				}
			}
		}

		for _, f := range failures {
			t.Errorf(
				"mismatched diagnostic at %s:\n"+
					"  want: [%s] %s\n"+
					"   got: [%s] %s",
				f.Expectation.Pos,
				f.Expectation.Severity,
				f.Expectation.Message,
				f.Diagnostic.Severity,
				f.Diagnostic.Message,
			)
		}

		for _, e := range missing {
			t.Errorf(
				"missing diagnostic at %s:\n"+
					"  want: [%s] %s",
				e.Pos,
				e.Severity,
				e.Message,
			)
		}

		for _, d := range unexpected {
			t.Errorf(
				"unexpected diagnostic at %s:\n"+
					"   got: [%s] %s",
				d.Begin,
				d.Severity,
				d.Message,
			)
		}
	})
}

func remove[S ~[]E, E comparable](s *S, e E) {
	for i, x := range *s {
		if x == e {
			*s = slices.Delete(*s, i, i+1)
			return
		}
	}
}
