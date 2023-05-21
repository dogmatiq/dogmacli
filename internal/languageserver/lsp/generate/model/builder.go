package model

import "reflect"

type builder struct {
	model     *Model
	parent    Node
	resolvers []func()
}

func build[Out Node](
	b *builder,
	fn func(Out),
) Out {
	out := newNode[Out]()

	p := b.parent
	out.setParent(p)

	b.parent = out
	fn(out)
	b.parent = p

	return out
}

func buildDef[In any, Out Def](
	b *builder,
	name string,
	in In,
	fn func(In, Out),
) {
	out := newNode[Out]()

	b.model.Defs[name] = out

	b.resolvers = append(
		b.resolvers,
		func() {
			p := b.parent
			out.setParent(p)

			b.parent = out
			fn(in, out)
			b.parent = p
		},
	)
}

func newNode[T Node]() T {
	var zero T
	return reflect.New(
		reflect.TypeOf(zero).Elem(),
	).Interface().(T)
}
