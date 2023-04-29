package identity

import (
	"go/ast"
	"go/constant"
	"strconv"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogmacli/internal/linter"
	"github.com/google/uuid"
	"golang.org/x/tools/go/ssa"
)

// Lint checks that all applications and handlers have a valid identity.
func Lint(ctx *linter.Context) {
	for _, e := range ctx.Entities.All {
		lint(ctx, e.ConfigureMethod)
	}
}

// lint checks that the identity of the given entity is valid.
func lint(ctx *linter.Context, c linter.ConfigureMethod) {
	visited := map[*ssa.BasicBlock]bool{}

	called := lintBlock(
		ctx,
		c,
		c.Implementation.Blocks[0],
		visited,
		0,
	)

	if !called {
		ctx.Error(
			c.Declaration,
			"Configure() must call %s.Identity() exactly once",
			c.Configurer.Name(),
		)
	}
}

// lintBlock looks for a call to the Identity() function in a specific block,
// and its successors.
//
// It returns true if Identity() is called at all, even if it is is not called
// on all execution paths.
func lintBlock(
	ctx *linter.Context,
	c linter.ConfigureMethod,
	block *ssa.BasicBlock,
	visited map[*ssa.BasicBlock]bool,
	priorCalls int,
) (called bool) {
	if called, ok := visited[block]; ok {
		return called
	}
	defer func() {
		visited[block] = called
	}()

	var calls []*ast.CallExpr

	for _, i := range block.Instrs {
		if call, ok := i.(*ssa.Call); ok {
			if c.IsConfigurerCall(call, "Identity") {
				called = true
				refs := *call.Referrers()
				expr := refs[0].(*ssa.DebugRef).Expr.(*ast.CallExpr)
				calls = append(calls, expr)

				lintName(ctx, expr.Args[0])
				lintKey(ctx, expr.Args[1])
			}
		}
	}

	// If there is more than one call in this block, all but the first is a
	// duplicate.
	//
	// If there are any prior calls at all (from parent blocks), then all of the
	// calls in this block are duplicates.
	if len(calls) > 0 {
		index := 0
		if priorCalls != 0 {
			index++
		}

		for _, call := range calls[index:] {
			ctx.Error(
				call,
				"%s.Identity() must be called exactly once",
				c.Configurer.Name(),
			)
		}
	}

	if len(block.Succs) == 0 {
		return called
	}

	if len(block.Succs) == 1 {
		nextCalled := lintBlock(
			ctx,
			c,
			block.Succs[0],
			visited,
			priorCalls+len(calls),
		)
		return called || nextCalled
	}

	thenBlock := block.Succs[0]
	thenCalled := lintBlock(
		ctx,
		c,
		thenBlock,
		visited,
		priorCalls+len(calls),
	)

	elseBlock := block.Succs[1]
	elseCalled := lintBlock(
		ctx,
		c,
		elseBlock,
		visited,
		priorCalls+len(calls),
	)

	if thenCalled != elseCalled {
		cond := block.Instrs[len(block.Instrs)-1].(*ssa.If).Cond
		refs := *cond.Referrers()
		expr := refs[0].(*ssa.DebugRef).Expr

		ctx.Error(
			expr,
			"this control-flow statement causes %s.Identity() to remain uncalled on some execution paths",
			c.Configurer.Name(),
		)
	}

	return called || thenCalled || elseCalled
}

func lintName(
	ctx *linter.Context,
	expr ast.Expr,
) {
	v := ctx.TypeInfo.Types[expr].Value
	if v == nil {
		return
	}

	name := constant.StringVal(v)

	placeholder := "0b24b57d-19f3-472b-bec1-70915974dadc"
	if _, err := configkit.NewIdentity(name, placeholder); err != nil {
		ctx.Error(
			expr,
			"%s",
			err.Error(),
		)
	}
}

func lintKey(
	ctx *linter.Context,
	expr ast.Expr,
) {
	v := ctx.TypeInfo.Types[expr].Value
	if v == nil {
		return
	}

	key := constant.StringVal(v)

	if _, err := configkit.NewIdentity("<placeholder>", key); err != nil {
		ctx.Error(
			expr,
			"%s",
			err.Error(),
		).SuggestChange(
			"generate a new UUID to use as the identity key",
			linter.Edit{
				Begin: expr.Pos(),
				End:   expr.End(),
				Text:  strconv.Quote(uuid.NewString()),
			},
		)
	}
}
