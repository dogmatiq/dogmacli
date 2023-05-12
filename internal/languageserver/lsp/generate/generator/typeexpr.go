package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/metamodel"
)

// typeExpr returns a Go type expression that represents the given meta-model
// type.
func typeExpr(t metamodel.Type) jen.Code {
	var v typeExprGenerator
	metamodel.VisitType(t, &v)
	return v.Code
}

// typeExprGenerator is an implementation of metamodel.TypeVisitor that
// generates code Go type expressions for meta-model types.
type typeExprGenerator struct {
	Code jen.Code
}

func (g *typeExprGenerator) VisitBool(t metamodel.BoolType) {
	g.Code = jen.Bool()
}

func (g *typeExprGenerator) VisitDecimal(t metamodel.DecimalType) {
	g.Code = jen.Float64()
}

func (g *typeExprGenerator) VisitString(t metamodel.StringType) {
	g.Code = jen.String()
}

func (g *typeExprGenerator) VisitInteger(t metamodel.IntegerType) {
	g.Code = jen.Int32()
}

func (g *typeExprGenerator) VisitUInteger(t metamodel.UIntegerType) {
	g.Code = jen.Uint32()
}

func (g *typeExprGenerator) VisitDocumentURI(t metamodel.DocumentURIType) {
	g.Code = jen.String()
}

func (g *typeExprGenerator) VisitURI(t metamodel.URIType) {
	g.Code = jen.String()
}

func (g *typeExprGenerator) VisitNull(t metamodel.NullType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitReference(t metamodel.ReferenceType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitArray(t metamodel.ArrayType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitMap(t metamodel.MapType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitLiteral(t metamodel.LiteralType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitStringLiteral(t metamodel.StringLiteralType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitAnd(t metamodel.AndType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitOr(t metamodel.OrType) {
	panic("not implemented")
}

func (g *typeExprGenerator) VisitTuple(t metamodel.TupleType) {
	panic("not implemented")
}
