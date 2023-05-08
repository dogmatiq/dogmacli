package identity

import (
	"go/ast"
	"go/constant"
	"strconv"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"github.com/dogmatiq/dogmacli/internal/linter/internal/lint"
	"github.com/google/uuid"
	"golang.org/x/tools/go/ssa"
)

// Lint checks that all applications and handlers have a valid identity.
func Lint(ctx *lint.Context) {
	visited := map[*ssa.BasicBlock]bool{}

	for _, e := range ctx.Entities {
		cm := e.ConfigureMethod

		called := lintBlock(
			ctx,
			cm,
			cm.Implementation.Blocks[0],
			visited,
			false,
		)

		if !called {
			ctx.Error(
				cm.Declaration,
				"Configure() must call %s.Identity()",
				cm.Configurer.Name(),
			)
		}
	}
}

// lintName checks that the name argument to Identity() is valid.
func lintName(ctx *lint.Context, expr ast.Expr) {
	v := ctx.Package.TypesInfo.Types[expr].Value
	if v == nil {
		return
	}

	name := constant.StringVal(v)

	if err := configkit.ValidateIdentityName(name); err != nil {
		ctx.Error(
			expr,
			"%s",
			err.Error(),
		)
	}
}

// lintKey checks that the key argument to Identity() is valid.
func lintKey(ctx *lint.Context, expr ast.Expr) {
	v := ctx.Package.TypesInfo.Types[expr].Value
	if v == nil {
		return
	}

	key := constant.StringVal(v)

	if err := configkit.ValidateIdentityKey(key); err != nil {
		ctx.Error(
			expr,
			"%s",
			err.Error(),
		).SuggestChange(
			"generate a new UUID to use as the identity key",
			diagnostic.Edit{
				Begin: expr.Pos(),
				End:   expr.End(),
				Text:  strconv.Quote(uuid.NewString()),
			},
		)
	}
}

// lintBlock looks for a call to the Identity() function in a specific block,
// and its successors.
//
// It returns true if Identity() is called at all, even if it is is not called
// on all execution paths.
func lintBlock(
	ctx *lint.Context,
	cm lint.ConfigureMethod,
	block *ssa.BasicBlock,
	visited map[*ssa.BasicBlock]bool,
	calledInPriorBlock bool,
) bool {
	calledInThisBlock, ok := visited[block]
	if ok {
		return calledInThisBlock
	}
	defer func() {
		visited[block] = calledInThisBlock
	}()

	var calls []*ast.CallExpr

	for _, i := range block.Instrs {
		if c, ok := i.(*ssa.Call); ok {
			if !cm.IsConfigurerCall(c, "Identity") {
				continue
			}

			refs := *c.Referrers()
			expr := refs[0].(*ssa.DebugRef).Expr.(*ast.CallExpr)
			calls = append(calls, expr)

			lintName(ctx, expr.Args[0])
			lintKey(ctx, expr.Args[1])

			if calledInThisBlock {
				ctx.Error(
					expr,
					"%s.Identity() must be called exactly once",
					cm.Configurer.Name(),
				)
			}

			calledInThisBlock = true
		}
	}

	if calledInThisBlock && calledInPriorBlock {
		ctx.Error(
			calls[0],
			"%s.Identity() has already been called on at least one execution path",
			cm.Configurer.Name(),
		)
	}

	if len(block.Succs) == 0 {
		return calledInThisBlock
	}

	calledInNextBlock := lintBlock(
		ctx,
		cm,
		block.Succs[0],
		visited,
		calledInThisBlock || calledInPriorBlock,
	)

	if len(block.Succs) == 1 {
		return calledInThisBlock || calledInNextBlock
	}

	calledInElseBlock := lintBlock(
		ctx,
		cm,
		block.Succs[1],
		visited,
		calledInThisBlock || calledInPriorBlock,
	)

	if calledInNextBlock != calledInElseBlock {
		cond := block.Instrs[len(block.Instrs)-1].(*ssa.If).Cond
		refs := *cond.Referrers()
		expr := refs[0].(*ssa.DebugRef).Expr

		ctx.Error(
			expr,
			"this control-flow statement causes %s.Identity() to remain uncalled on some execution paths",
			cm.Configurer.Name(),
		)
	}

	return calledInThisBlock || calledInNextBlock || calledInElseBlock
}
