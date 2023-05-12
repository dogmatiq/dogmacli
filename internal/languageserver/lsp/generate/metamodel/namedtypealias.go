package metamodel

// TypeAlias defines a named type alias.
type TypeAlias struct {
	Name          string
	Documentation string
	Type          Type
}

func populateTypeAlias(
	named map[string]NamedType,
	j typeAliasJSON,
) {
	t := named[j.Name].(*TypeAlias)

	t.Name = j.Name
	t.Documentation = j.Documentation
	t.Type = newType(named, j.Type)
}

// TypeName returns the name type's name.
func (t *TypeAlias) TypeName() string {
	return t.Name
}

type typeAliasJSON struct {
	Name          string    `json:"name"`
	Documentation string    `json:"documentation"`
	Type          *typeJSON `json:"type"`
}
