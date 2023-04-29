package linter

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

type (
	// Entity is an application or handler.
	Entity struct {
		ConfigureMethod ConfigureMethod
	}

	// Application contains information about a type that implements
	// [dogma.Application].
	Application = Entity

	// Handler contains information about a type that implements one of Dogma's
	// handler interfaces.
	Handler = Entity

	// Aggregate contains information about a type that implements
	// [dogma.AggregateMessageHandler].
	Aggregate = Handler

	// Process contains information about a type that implements
	// [dogma.ProcessMessageHandler].
	Process = Handler

	// Integration contains information about a type that implements
	// [dogma.IntegrationMessageHandler].
	Integration = Handler

	// Projection contains information about a type that implements
	// [dogma.ProjectionMessageHandler].
	Projection = Handler
)

// ConfigureMethod describes the Configure() method of an application or
// handler.
type ConfigureMethod struct {
	Declaration    *ast.FuncDecl
	Signature      *types.Func
	Implementation *ssa.Function
	Configurer     *ssa.Parameter
}

// IsConfigurerCall returns true if c represents a call to a specific method on
// the configurer.
func (m ConfigureMethod) IsConfigurerCall(c *ssa.Call, method string) bool {
	panic("not implemented")
}
