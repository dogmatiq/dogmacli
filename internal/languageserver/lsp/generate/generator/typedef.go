package generator

import "github.com/dave/jennifer/jen"

// // typeName returns the Go type name to use for the given definition.
// func typeName(d model.TypeDef) string {
// 	return normalized(d.Name())
// }

type typeDefGen struct {
	*jen.File
}
