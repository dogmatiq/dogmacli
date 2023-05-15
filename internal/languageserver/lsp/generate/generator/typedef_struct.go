package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g typeDefGen) Struct(d model.Struct) {
	documentation(g, d.Documentation)
	g.
		Type().
		Id(exported(d.TypeName)).
		StructFunc(func(g *jen.Group) {
			for _, t := range d.EmbeddedTypes {
				g.Add(typeExpr(t))
			}

			for _, p := range d.Properties {
				if _, ok := p.Type.(model.StringLit); !ok {
					documentation(g, p.Documentation)

					t := typeExpr(p.Type)
					if p.Optional && useOptionalType(p.Type) {
						t = jen.Id("Optional").Types(t)
					}

					g.
						Id(exported(p.Name)).
						Add(t)
				}
			}
		})

	g.Line()
	g.structValidateMethod(d)
}

func (g typeDefGen) structValidateMethod(d model.Struct) {
	g.Comment("Validate returns an error if x is invalid.")
	g.
		Func().
		Params(
			jen.Id("x").Id(exported(d.TypeName)),
		).
		Id("Validate").
		Params().
		Params(
			jen.Error(),
		).
		BlockFunc(func(g *jen.Group) {
			g.Panic(jen.Lit("not implemented"))
		})
}

// useOptionalType returns the Go type expression that refers to t.
func useOptionalType(t model.Type) bool {
	return model.ApplyTypeTransform[bool](
		t,
		useOptionalTypeX{},
	)
}

type useOptionalTypeX struct{}

func (useOptionalTypeX) Bool() bool                       { return false }
func (useOptionalTypeX) Decimal() bool                    { return false }
func (useOptionalTypeX) String() bool                     { return false }
func (useOptionalTypeX) Integer() bool                    { return false }
func (useOptionalTypeX) UInteger() bool                   { return false }
func (useOptionalTypeX) DocumentURI() bool                { return false }
func (useOptionalTypeX) URI() bool                        { return false }
func (useOptionalTypeX) Null() bool                       { return false }
func (useOptionalTypeX) Reference(t model.Reference) bool { return true }
func (useOptionalTypeX) Array(t model.Array) bool         { return false }
func (useOptionalTypeX) Map(t model.Map) bool             { return false }
func (useOptionalTypeX) And(t model.And) bool             { return true }
func (useOptionalTypeX) Or(t model.Or) bool               { return true }
func (useOptionalTypeX) Tuple(t model.Tuple) bool         { return true }
func (useOptionalTypeX) StructLit(t model.StructLit) bool { return true }
func (useOptionalTypeX) StringLit(t model.StringLit) bool { return false }
