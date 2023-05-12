package metamodel

// Structure defines a named structure data type.
type Structure struct {
	Name          string
	Documentation string
	EmbeddedTypes []ReferenceType
	Properties    []StructureProperty
}

func populateStructure(
	named map[string]NamedType,
	j structureJSON,
) {
	t := named[j.Name].(*Structure)

	t.Name = j.Name
	t.Documentation = j.Documentation
	t.Properties = newStructureProperties(named, j.Properties)

	for _, embed := range j.Extends {
		t.EmbeddedTypes = append(
			t.EmbeddedTypes,
			ReferenceType{
				Target: named[embed.Name],
			},
		)
	}

	for _, embed := range j.Mixins {
		t.EmbeddedTypes = append(
			t.EmbeddedTypes,
			ReferenceType{
				Target: named[embed.Name],
			},
		)
	}
}

// TypeName returns the name type's name.
func (s *Structure) TypeName() string {
	return s.Name
}

// StructureProperty is property within a structure.
type StructureProperty struct {
	Name          string
	Documentation string
	Optional      bool
	Type          Type
}

func newStructureProperties(
	named map[string]NamedType,
	j []structurePropertyJSON,
) []StructureProperty {
	var properties []StructureProperty

	for _, p := range j {
		properties = append(
			properties,
			StructureProperty{
				Name:          p.Name,
				Documentation: p.Documentation,
				Optional:      p.Optional,
				Type:          newType(named, p.Type),
			},
		)
	}

	return properties
}

type structureJSON struct {
	Name          string                  `json:"name"`
	Documentation string                  `json:"documentation"`
	Properties    []structurePropertyJSON `json:"properties"`
	Extends       []*typeJSON             `json:"extends"`
	Mixins        []*typeJSON             `json:"mixins"`
}

type structurePropertyJSON struct {
	Name          string    `json:"name"`
	Documentation string    `json:"documentation"`
	Optional      bool      `json:"optional"`
	Type          *typeJSON `json:"type"`
}
