package lsp

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

const (
	// dogmaPkgPath is the full path of dogma package.
	dogmaPkgPath = "github.com/dogmatiq/dogma"
)

// dogmaPackage encapsulates information about the dogma package and the types
// within it.
type dogmaPackage struct {
	Package                   *types.Package
	Application               *types.Interface
	AggregateMessageHandler   *types.Interface
	ProcessMessageHandler     *types.Interface
	ProjectionMessageHandler  *types.Interface
	IntegrationMessageHandler *types.Interface
}

// lookupDogmaPackage returns information about the dogma package.
//
// It returns false if Dogma has not been imported.
func lookupDogmaPackage(prog *ssa.Program) (dogmaPackage, bool) {
	dogmaPkg := prog.ImportedPackage(dogmaPkgPath)
	if dogmaPkg == nil {
		return dogmaPackage{}, false
	}

	scope := dogmaPkg.Pkg.Scope()

	return dogmaPackage{
		Package:                   dogmaPkg.Pkg,
		Application:               scope.Lookup("Application").Type().Underlying().(*types.Interface),
		AggregateMessageHandler:   scope.Lookup("AggregateMessageHandler").Type().Underlying().(*types.Interface),
		ProcessMessageHandler:     scope.Lookup("ProcessMessageHandler").Type().Underlying().(*types.Interface),
		ProjectionMessageHandler:  scope.Lookup("ProjectionMessageHandler").Type().Underlying().(*types.Interface),
		IntegrationMessageHandler: scope.Lookup("IntegrationMessageHandler").Type().Underlying().(*types.Interface),
	}, true
}
