package jenx

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

// Litf generates a string literal.
func Litf(
	format string,
	args ...any,
) *jen.Statement {
	return jen.Lit(
		fmt.Sprintf(format, args...),
	)
}

// Errorf generates a call to fmt.Errorf().
func Errorf(
	format string,
	args ...jen.Code,
) *jen.Statement {
	return jen.
		Qual("fmt", "Errorf").
		CallFunc(func(grp *jen.Group) {
			grp.Lit(format)
			for _, arg := range args {
				grp.Add(arg)
			}
		})
}

// Sprintf generates a call to fmt.Sprintf().
func Sprintf(
	format string,
	args ...jen.Code,
) *jen.Statement {
	return jen.
		Qual("fmt", "Sprintf").
		CallFunc(func(grp *jen.Group) {
			grp.Lit(format)
			for _, arg := range args {
				grp.Add(arg)
			}
		})
}
