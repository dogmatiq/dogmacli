package lint

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

// Context encapsulates the information available to a linter, and provides a
// way for the linter to report diagnostics.
type Context struct {
	Package    *packages.Package
	SSAProgram *ssa.Program
	SSAPackage *ssa.Package

	Dogma        *Dogma
	Entities     []Entity
	Applications []Application
	Handlers     []Handler
	Aggregates   []Aggregate
	Processes    []Process
	Integrations []Integration
	Projections  []Projection

	Diagnostics []*diagnostic.Diagnostic
}

// NewContext creates a new linter context for the linting the given package.
func NewContext(
	pkg *packages.Package,
	prog *ssa.Program,
	spkg *ssa.Package,
) *Context {
	ctx := &Context{
		Package:    pkg,
		SSAProgram: prog,
		SSAPackage: spkg,
	}

	if d := ctx.SSAProgram.ImportedPackage(dogmaPkgPath); d != nil {
		iface := func(n string) *types.Interface {
			return d.
				Pkg.
				Scope().
				Lookup(n).
				Type().
				Underlying().(*types.Interface)
		}

		ctx.Dogma = &Dogma{
			Package:                   d.Pkg,
			Application:               iface("Application"),
			AggregateMessageHandler:   iface("AggregateMessageHandler"),
			ProcessMessageHandler:     iface("ProcessMessageHandler"),
			ProjectionMessageHandler:  iface("ProjectionMessageHandler"),
			IntegrationMessageHandler: iface("IntegrationMessageHandler"),
		}

		buildEntities(ctx)
	}

	return ctx
}

// Report reports a diagnostic about the given AST node.
func (c *Context) Report(
	s diagnostic.Severity,
	n ast.Node,
	format string, args ...any,
) *diagnostic.Diagnostic {
	d := &diagnostic.Diagnostic{
		Begin:   c.Package.Fset.Position(n.Pos()),
		End:     c.Package.Fset.Position(n.End()),
		Message: fmt.Sprintf(format, args...),
	}

	c.Diagnostics = append(c.Diagnostics, d)

	return d
}

// Error reports an error about the given AST node.
func (c *Context) Error(
	n ast.Node,
	format string, args ...any,
) *diagnostic.Diagnostic {
	return c.Report(diagnostic.Error, n, format, args...)
}

// Warning reports a warning about the given AST node.
func (c *Context) Warning(
	n ast.Node,
	format string, args ...any,
) *diagnostic.Diagnostic {
	return c.Report(diagnostic.Warning, n, format, args...)
}

// Improvement reports an improvement that can be made to the given AST node.
func (c *Context) Improvement(
	n ast.Node,
	format string, args ...any,
) *diagnostic.Diagnostic {
	return c.Report(diagnostic.Improvement, n, format, args...)
}
