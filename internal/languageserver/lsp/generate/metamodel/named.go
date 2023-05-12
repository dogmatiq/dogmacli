package metamodel

// NamedType is an interface for definitions of named types.
type NamedType interface {
	TypeName() string
	acceptVisitor(NamedTypeVisitor)
}

// VisitNamedType dispatches to the method on v that corresponds to t's type.
func VisitNamedType(t NamedType, v NamedTypeVisitor) {
	t.acceptVisitor(v)
}

// NamedTypeVisitor provides named-type-specific logic.
type NamedTypeVisitor interface {
	VisitEnumeration(Enumeration)
	VisitStructure(Structure)
	VisitTypeAlias(TypeAlias)
}

func (t *Enumeration) acceptVisitor(v NamedTypeVisitor) { v.VisitEnumeration(*t) }
func (t *Structure) acceptVisitor(v NamedTypeVisitor)   { v.VisitStructure(*t) }
func (t *TypeAlias) acceptVisitor(v NamedTypeVisitor)   { v.VisitTypeAlias(*t) }
