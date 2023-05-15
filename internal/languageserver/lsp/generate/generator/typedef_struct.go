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
	return model.TransformType[bool](
		t,
		useOptX{},
	)
}

type useOptX struct{}

func (x useOptX) Bool() bool                       { return false }
func (x useOptX) Decimal() bool                    { return false }
func (x useOptX) String() bool                     { return false }
func (x useOptX) Integer() bool                    { return false }
func (x useOptX) UInteger() bool                   { return false }
func (x useOptX) DocumentURI() bool                { return false }
func (x useOptX) URI() bool                        { return false }
func (x useOptX) Null() bool                       { return false }
func (x useOptX) Array(t model.Array) bool         { return false }
func (x useOptX) Map(t model.Map) bool             { return false }
func (x useOptX) And(t model.And) bool             { return true }
func (x useOptX) Or(t model.Or) bool               { return true }
func (x useOptX) Tuple(t model.Tuple) bool         { return true }
func (x useOptX) StructLit(t model.StructLit) bool { return true }
func (x useOptX) StringLit(t model.StringLit) bool { return false }
func (x useOptX) Reference(t model.Reference) bool { return model.TransformTypeDef[bool](t.Target, x) }
func (x useOptX) Alias(d model.Alias) bool         { return model.TransformType[bool](d.Type, x) }
func (x useOptX) Enum(d model.Enum) bool           { return true }
func (x useOptX) Struct(d model.Struct) bool       { return true }
