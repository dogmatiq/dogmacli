package model

// TypeDef is a definition of a named type.
type TypeDef interface {
	TypeName() string
	isTypeDef()
}

type typeDef struct {
	node
	name string

	Documentation Documentation
}

func (d typeDef) TypeName() string {
	return d.name
}

func (typeDef) isTypeDef() {}
