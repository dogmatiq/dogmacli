package lint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

type (
	// Entity is an application or handler.
	Entity struct {
		Type            types.Type
		ConfigureMethod ConfigureMethod
	}

	// Application contains information about a type that implements
	// [dogma.Application].
	Application struct{ Entity }

	// Handler contains information about a type that implements one of Dogma's
	// handler interfaces.
	Handler = struct{ Entity }

	// Aggregate contains information about a type that implements
	// [dogma.AggregateMessageHandler].
	Aggregate struct{ Handler }

	// Process contains information about a type that implements
	// [dogma.ProcessMessageHandler].
	Process struct{ Handler }

	// Integration contains information about a type that implements
	// [dogma.IntegrationMessageHandler].
	Integration struct{ Handler }

	// Projection contains information about a type that implements
	// [dogma.ProjectionMessageHandler].
	Projection struct{ Handler }
)

// ConfigureMethod describes the Configure() method of an application or
// handler.
type ConfigureMethod struct {
	Declaration    *ast.FuncDecl
	Implementation *ssa.Function
	Configurer     *ssa.Parameter
}

// IsConfigurerCall returns true if c is a call to a specific method of the
// entity's configurer.
func (m *ConfigureMethod) IsConfigurerCall(c *ssa.Call, method string) bool {
	return c.Call.Value == m.Configurer && c.Common().Method.Name() == method
}

// buildEntities populates the entities fields of ctx.
func buildEntities(ctx *Context) {
	// If the package does not import dogma, then there is nothing to do.
	if ctx.Dogma == nil {
		return
	}

	for _, m := range ctx.SSAPackage.Members {
		if m, ok := m.(*ssa.Type); ok {
			t := m.Type()

			// TODO: skip type aliases

			if t, ok := implements(t, ctx.Dogma.Application); ok {
				buildApplication(ctx, t)
			} else if t, ok := implements(t, ctx.Dogma.AggregateMessageHandler); ok {
				buildAggregate(ctx, t)
			} else if t, ok := implements(t, ctx.Dogma.ProcessMessageHandler); ok {
				buildProcess(ctx, t)
			} else if t, ok := implements(t, ctx.Dogma.IntegrationMessageHandler); ok {
				buildIntegration(ctx, t)
			} else if t, ok := implements(t, ctx.Dogma.ProjectionMessageHandler); ok {
				buildProjection(ctx, t)
			}
		}
	}
}

// implements returns true if t implements i.
//
// It returns the actual type that implements the interface, either t or *t.
func implements(t types.Type, i *types.Interface) (types.Type, bool) {
	// The sequence of the if-blocks below is important as a type
	// implements an interface only if the methods in the interface's
	// method set have non-pointer receivers. Hence the implementation
	// check for the "raw" (non-pointer) type is made first.
	//
	// A pointer to the type, on the other hand, implements the
	// interface regardless of whether pointer receivers are used or
	// not.
	if types.Implements(t, i) {
		return t, true
	}

	return t, types.Implements(t, i)
}

func buildEntity(ctx *Context, t types.Type) Entity {
	fn := ctx.SSAProgram.LookupMethod(t, ctx.SSAPackage.Pkg, "Configure")

	e := Entity{
		Type: t,
		ConfigureMethod: ConfigureMethod{
			Declaration:    fn.Syntax().(*ast.FuncDecl),
			Implementation: fn,
			Configurer:     fn.Params[1],
		},
	}

	ctx.Entities = append(ctx.Entities, e)

	return e
}

func buildApplication(ctx *Context, t types.Type) {
	ctx.Applications = append(
		ctx.Applications,
		Application{
			Entity: buildEntity(ctx, t),
		},
	)
}

func buildHandler(ctx *Context, t types.Type) Handler {
	h := Handler{
		Entity: buildEntity(ctx, t),
	}

	ctx.Handlers = append(ctx.Handlers, h)

	return h
}

func buildAggregate(ctx *Context, t types.Type) {
	ctx.Aggregates = append(
		ctx.Aggregates,
		Aggregate{
			Handler: buildHandler(ctx, t),
		},
	)
}

func buildProcess(ctx *Context, t types.Type) {
	ctx.Processes = append(
		ctx.Processes,
		Process{
			Handler: buildHandler(ctx, t),
		},
	)
}

func buildIntegration(ctx *Context, t types.Type) {
	ctx.Integrations = append(
		ctx.Integrations,
		Integration{
			Handler: buildHandler(ctx, t),
		},
	)
}

func buildProjection(ctx *Context, t types.Type) {
	ctx.Projections = append(
		ctx.Projections,
		Projection{
			Handler: buildHandler(ctx, t),
		},
	)
}
