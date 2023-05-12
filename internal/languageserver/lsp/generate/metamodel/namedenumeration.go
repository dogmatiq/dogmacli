package metamodel

import "math"

// Enumeration defines a named enumeration data type.
type Enumeration struct {
	Name          string
	Documentation string
	Type          Type
	Members       []EnumerationMember
}

func populateEnumeration(
	named map[string]NamedType,
	j enumerationJSON,
) {
	t := named[j.Name].(*Enumeration)

	t.Name = j.Name
	t.Documentation = j.Documentation
	t.Type = newType(named, j.Type)

	for _, m := range j.Members {
		v := m.Value
		if f, ok := v.(float64); ok {
			if math.Mod(f, 1) == 0 {
				v = int(f)
			}
		}

		t.Members = append(
			t.Members,
			EnumerationMember{
				Name:          m.Name,
				Documentation: m.Documentation,
				Value:         v,
			},
		)
	}
}

// TypeName returns the name type's name.
func (t *Enumeration) TypeName() string {
	return t.Name
}

// EnumerationMember is a member of an enumeration.
type EnumerationMember struct {
	Name          string
	Documentation string
	Value         any
}

type enumerationJSON struct {
	Name          string                  `json:"name"`
	Documentation string                  `json:"documentation"`
	Type          *typeJSON               `json:"type"`
	Members       []enumerationMemberJSON `json:"values"`
}

type enumerationMemberJSON struct {
	Name          string `json:"name"`
	Documentation string `json:"documentation"`
	Value         any    `json:"value"`
}
