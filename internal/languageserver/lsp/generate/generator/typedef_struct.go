package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

func (g *typeDef) Struct(d model.Struct) {
	documentation(g.File, d.Documentation)
	g.File.
		Type().
		Id(identifier(d.TypeName)).
		StructFunc(func(grp *jen.Group) {
			for _, t := range d.EmbeddedTypes {
				grp.Add(g.typeExpr(t))
			}

			for _, p := range d.Properties {
				g.structProperty(grp, p)
			}
		})

	// g.File.Line()
	// g.structValidateMethod(d)
}

func (g *Generator) structProperty(grp *jen.Group, p model.Property) {
	if _, ok := p.Type.(model.StringLit); ok {
		return
	}

	g.pushName(p.Name)
	defer g.popName()

	documentation(grp, p.Documentation)

	i := g.typeInfo(p.Type)
	t := g.typeExpr(p.Type)
	if p.Optional && i.UseOptional {
		t = jen.Id("Optional").Types(t)
	}

	grp.
		Id(identifier(p.Name)).
		Add(t)
}

func (g *Generator) structValidateMethod(d model.Struct) {
	g.File.
		Comment("Validate returns an error if x is invalid.").
		Func().
		Params(
			jen.Id("x").Id(identifier(d.TypeName)),
		).
		Id("Validate").
		Params().
		Params(
			jen.Error(),
		).
		BlockFunc(func(grp *jen.Group) {
			grp.Panic(jen.Lit("not implemented"))
		})
}
