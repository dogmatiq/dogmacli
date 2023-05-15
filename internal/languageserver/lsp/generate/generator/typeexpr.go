package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

// typeExpr returns the Go type expression that refers to t.
func typeExpr(t model.Type) *jen.Statement {
	return model.ApplyTypeTransform[*jen.Statement](
		t,
		typeExprX{},
	)
}

type typeExprX struct{}

func (typeExprX) Bool() *jen.Statement                       { return jen.Bool() }
func (typeExprX) Decimal() *jen.Statement                    { return jen.Float64() }
func (typeExprX) String() *jen.Statement                     { return jen.String() }
func (typeExprX) Integer() *jen.Statement                    { return jen.Int32() }
func (typeExprX) UInteger() *jen.Statement                   { return jen.Uint32() }
func (typeExprX) DocumentURI() *jen.Statement                { return jen.Id("DocumentURI") }
func (typeExprX) URI() *jen.Statement                        { return jen.Id("URI") }
func (typeExprX) Null() *jen.Statement                       { return jen.Id("Null") }
func (typeExprX) Reference(t model.Reference) *jen.Statement { return jen.Id(t.Target.Name()) }

func (typeExprX) Array(t model.Array) *jen.Statement {
	return jen.
		Index().
		Add(typeExpr(t.Element))
}

func (typeExprX) Map(t model.Map) *jen.Statement {
	return jen.
		Map(typeExpr(t.Key)).
		Add(typeExpr(t.Value))
}

func (typeExprX) And(t model.And) *jen.Statement             { return jen.Struct() } // TODO
func (typeExprX) Or(t model.Or) *jen.Statement               { return jen.Struct() } // TODO
func (typeExprX) Tuple(t model.Tuple) *jen.Statement         { return jen.Struct() } // TODO
func (typeExprX) StructLit(t model.StructLit) *jen.Statement { return jen.Struct() } // TODO
func (typeExprX) StringLit(t model.StringLit) *jen.Statement { return jen.Id(unexported(t.Value)) }
