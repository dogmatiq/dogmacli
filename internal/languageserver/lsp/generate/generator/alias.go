package generator

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"

func (g *generator) VisitAlias(n *model.Alias) {
	// Any alias for an anonymous type is ignored. Instead, the underlying type
	// assumes the name of the alias.
	if n.UnderlyingType.IsAnonymous() {
		return
	}

	name := nameOf(n)
	underlying := nameOf(n.UnderlyingType)

	documentation(
		g,
		n.Documentation(),
		"%s is an alias for %s.",
		name,
		underlying,
	)

	g.
		Type().
		Id(name).
		Id(underlying)
}
