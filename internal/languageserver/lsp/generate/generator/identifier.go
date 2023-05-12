package generator

import "strings"

// exported returns a normalized exported identifier containing the given parts.
func exported(parts ...string) string {
	var id string

	for _, p := range parts {
		p = strings.Title(p)
		id += p
	}

	return id
}
