package generator

import "strings"

// exported returns a normalized exported identifier containing the given parts.
func exported(parts ...string) string {
	var id string

	for _, p := range parts {
		p = strings.Title(p)
		id += p
	}

	if x, ok := strings.CutSuffix(id, "Id"); ok {
		id = x + "ID"
	}
	if x, ok := strings.CutSuffix(id, "Ids"); ok {
		id = x + "IDs"
	}
	if x, ok := strings.CutSuffix(id, "Uri"); ok {
		id = x + "URI"
	}

	return id
}

// unexported returns a normalized unexported identifier containing the given
// parts.
func unexported(parts ...string) string {
	id := exported(parts...)
	return strings.ToLower(id[:1]) + id[1:]
}

// identifier returns a identifier identifier containing the given parts.
func identifier(parts ...string) string {
	if parts[0][0] == '_' {
		parts[0] = parts[0][1:]
		return unexported(parts...)
	}
	return exported(parts...)
}
