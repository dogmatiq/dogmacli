package generator

import "strings"

var initialisms = []string{
	"ID",
	"URI",
}

// exported returns a normalized exported identifier containing the given parts.
func exported(parts ...string) string {
	var id string

	for _, p := range parts {
		p = strings.Title(p)
		id += p
	}

	for _, init := range initialisms {
		if strings.HasSuffix(
			strings.ToLower(id),
			strings.ToLower(init),
		) {
			id = id[:len(id)-len(init)] + init
		}
	}

	return id
}

// unexported returns a normalized unexported identifier containing the given
// parts.
func unexported(parts ...string) string {
	id := exported(parts...)
	return strings.ToLower(id[:1]) + id[1:]
}

// normalized returns a normalized identifier containing the given parts.
func normalized(parts ...string) string {
	if parts[0][0] == '_' {
		parts[0] = parts[0][1:]
		return unexported(parts...)
	}
	return exported(parts...)
}
