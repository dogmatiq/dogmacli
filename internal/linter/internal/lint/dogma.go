package lint

import (
	"go/types"
)

const (
	// dogmaPkgPath is the full path of dogma package.
	dogmaPkgPath = "github.com/dogmatiq/dogma"
)

// Dogma encapsulates information about the dogma package and the types within
// it.
type Dogma struct {
	Package *types.Package

	Application               *types.Interface
	AggregateMessageHandler   *types.Interface
	ProcessMessageHandler     *types.Interface
	ProjectionMessageHandler  *types.Interface
	IntegrationMessageHandler *types.Interface
}

// buildDogmaTypes populates the ctx.Dogma field if the dogma package is
// imported by ctx.Package.
func buildDogmaTypes(ctx *Context) {
	pkg := ctx.SSAProgram.ImportedPackage(dogmaPkgPath)
	if pkg == nil {
		return
	}

	iface := func(n string) *types.Interface {
		return pkg.
			Pkg.
			Scope().
			Lookup(n).
			Type().
			Underlying().(*types.Interface)
	}

	ctx.Dogma = &Dogma{
		Package:                   pkg.Pkg,
		Application:               iface("Application"),
		AggregateMessageHandler:   iface("AggregateMessageHandler"),
		ProcessMessageHandler:     iface("ProcessMessageHandler"),
		ProjectionMessageHandler:  iface("ProjectionMessageHandler"),
		IntegrationMessageHandler: iface("IntegrationMessageHandler"),
	}
}
