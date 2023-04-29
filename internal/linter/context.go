package linter

import "go/types"

// Context encapsulates the information available to a linter, and provides a
// way for the linter to report diagnostics.
type Context struct {
	Reporter

	Entities struct {
		All          []*Entity
		Applications []*Application
		Handlers     []*Handler
		Aggregates   []*Aggregate
		Processes    []*Process
		Integrations []*Integration
		Projections  []*Projection
	}

	// Package    *types.Package
	// Files      []*ast.File
	// FileSet    *token.FileSet
	TypeInfo *types.Info
	// SSAPackage *ssa.Package
}

// // buildSSA builds the SSA representation of the package under anaylsis.
// func buildSSA(ctx *Context) *ssa.Package {
// 	var (
// 		prog    = ssa.NewProgram(ctx.FileSet, 0)
// 		visited = map[*types.Package]struct{}{}
// 		imports func(pkgs *types.Package)
// 	)

// 	imports = func(pkg *types.Package) {
// 		for _, p := range pkg.Imports() {
// 			if _, ok := visited[p]; ok {
// 				return
// 			}
// 			visited[p] = struct{}{}
// 			prog.CreatePackage(p, nil, nil, true)
// 			imports(p)
// 		}
// 	}

// 	imports(ctx.Package)

// 	pkg := prog.CreatePackage(ctx.Package, ctx.Files, ctx.TypeInfo, false)
// 	pkg.SetDebugMode(true)
// 	pkg.Build()

// 	return pkg
// }

// // isConfigureImplementation returns true if fn is a concrete method with a
// // signature like:
// //
// //	func (T) Configure(dogma.[XXX]Configurer)
// func isConfigureImplementation(
// 	ctx *rules.Context,
// 	decl *ast.FuncDecl,
// ) bool {
// 	if decl.Recv == nil {
// 		// This is a function, not a method.
// 		return false
// 	}

// 	if decl.Name.Name != "Configure" {
// 		// This method is not named Configure().
// 		return false
// 	}

// 	if decl.Type.Params.NumFields() != 1 {
// 		// This function does not accept exactly one parameter (not including
// 		// the receiver).
// 		return false
// 	}

// 	param := decl.Type.Params.List[0]

// 	nt, ok := ctx.TypeInfo.TypeOf(param.Type).(*types.Named)
// 	if !ok {
// 		// The parameter does not have a named type (can't be dogma.Something).
// 		return false
// 	}

// 	if nt.Obj().Pkg().Path() != dogmatypes.PkgPath {
// 		// The parameter type is not in the Dogma package.
// 		return false
// 	}

// 	if !strings.HasSuffix(nt.Obj().Name(), "Configurer") {
// 		// The parameter type is not one of the Dogma configurer types.
// 		return false
// 	}

// 	return true
// }
