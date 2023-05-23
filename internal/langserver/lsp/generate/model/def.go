package model

// Def is the definition of a root-level named entity within the model.
type Def interface {
	Node

	Name() string
	Documentation() Documentation
}

// defNode provides implementation common to all types that implement Def.
type defNode struct {
	node
	name string
	docs Documentation
}

func (d defNode) Name() string {
	return d.name
}

func (d defNode) Documentation() Documentation {
	return d.docs
}
