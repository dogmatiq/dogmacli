package model

// TypeDef is a definition of a named type.
type TypeDef interface {
	Name() string
	accept(TypeDefVisitor)
}

// TypeDefVisitor provides logic specific to each TypeDef implementation.
type TypeDefVisitor interface {
	Alias(*Alias)
	Enum(*Enum)
	Struct(*Struct)
}

// VisitTypeDef dispatches to v based on the concrete type of d.
func VisitTypeDef(d TypeDef, v TypeDefVisitor) {
	d.accept(v)
}

// TypeDefTransform produces a value of type T from a TypeDef.
type TypeDefTransform[T any] interface {
	Alias(*Alias) T
	Enum(*Enum) T
	Struct(*Struct) T
}

// TypeDefTo transforms d to a value of type T using x.
func TypeDefTo[T any](
	d TypeDef,
	x TypeDefTransform[T],
) T {
	v := &typeDefX[T]{X: x}
	VisitTypeDef(d, v)
	return v.V
}

type typeDefX[T any] struct {
	X TypeDefTransform[T]
	V T
}

// typeDef returns the type definition with the given name.
func (b *builder) typeDef(name string) TypeDef {
	if t, ok := b.aliases[name]; ok {
		return t
	}

	if t, ok := b.structs[name]; ok {
		return t
	}

	if t, ok := b.enums[name]; ok {
		return t
	}

	panic("unknown type: " + name)
}
