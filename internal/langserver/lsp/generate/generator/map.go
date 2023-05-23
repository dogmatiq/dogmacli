package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model"
)

func (g *generator) VisitMap(n *model.Map) {
	name := nameOf(n)
	key := nameOf(n.KeyType)
	value := nameOf(n.ValueType)

	g.
		Commentf("%s is a map of %s to %s.", name, key, value).
		Line().
		Type().
		Id(name).
		Map(jen.Id(key)).
		Id(value)
}
