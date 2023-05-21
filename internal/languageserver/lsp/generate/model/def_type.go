package model

// TypeDef is a definition of a named type.
type TypeDef interface {
	Def

	isTypeDef()
}

// typeDefNode provides implementation common to all types that implement
// TypeDef.
type typeDefNode struct {
	defNode
}

func (typeDefNode) isTypeDef() {}
