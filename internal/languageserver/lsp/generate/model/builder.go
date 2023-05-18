package model

import "reflect"

type builder struct {
	stack []Node

	aliases map[string]*Alias
	enums   map[string]*Enum
	structs map[string]*Struct
}

func build[T Node](
	b *builder,
	fn func(n T),
) T {
	var node T
	node = reflect.New(
		reflect.TypeOf(node).Elem(),
	).Interface().(T)

	if len(b.stack) > 0 {
		p := b.stack[len(b.stack)-1]
		node.setParent(p)
	}

	b.stack = append(b.stack, node)
	fn(node)
	b.stack = b.stack[:len(b.stack)-1]

	return node
}
