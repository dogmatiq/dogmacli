package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

// typeExpr returns the Go type expression that refers to t.
func (g *Generator) typeExpr(t model.Type) *jen.Statement {
	return model.TypeTo[*jen.Statement](
		t,
		&typeExpr{g},
	)
}

type typeExpr struct{ *Generator }

func (g *typeExpr) Bool() *jen.Statement        { return jen.Bool() }
func (g *typeExpr) Decimal() *jen.Statement     { return jen.Float64() }
func (g *typeExpr) String() *jen.Statement      { return jen.String() }
func (g *typeExpr) Integer() *jen.Statement     { return jen.Int32() }
func (g *typeExpr) UInteger() *jen.Statement    { return jen.Uint32() }
func (g *typeExpr) DocumentURI() *jen.Statement { return jen.Id("DocumentURI") }
func (g *typeExpr) URI() *jen.Statement         { return jen.Id("URI") }
func (g *typeExpr) Null() *jen.Statement        { return jen.Id("Null") }

func (g *typeExpr) Reference(t model.Reference) *jen.Statement {
	return jen.Id(identifier(t.Target.Name()))
}

func (g *typeExpr) Array(t model.Array) *jen.Statement {
	return jen.
		Index().
		Add(g.typeExpr(t.Element))
}

func (g *typeExpr) Map(t model.Map) *jen.Statement {
	return jen.
		Map(g.typeExpr(t.Key)).
		Add(g.typeExpr(t.Value))
}

func (g *typeExpr) And(t model.And) *jen.Statement             { return g.reify(t) }
func (g *typeExpr) Or(t model.Or) *jen.Statement               { return g.reify(t) }
func (g *typeExpr) Tuple(t model.Tuple) *jen.Statement         { return g.reify(t) }
func (g *typeExpr) StructLit(t model.StructLit) *jen.Statement { return g.reify(t) }

func (g *typeExpr) StringLit(t model.StringLit) *jen.Statement {
	panic("string literals do not have a type representation")
}
