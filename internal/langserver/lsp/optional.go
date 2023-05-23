package lsp

// Optional is an optional value of type T.
type Optional[T any] struct {
	value *T
}

// Get returns the value, or false if the optional is empty.
func (o Optional[T]) Get() (T, bool) {
	if o.value != nil {
		return *o.value, true
	}

	var zero T
	return zero, false
}

// With returns an optional value that contains v.
func With[T any](v T) Optional[T] {
	return Optional[T]{&v}
}

// Without returns an empty optional value.
func Without[T any]() Optional[T] {
	return Optional[T]{}
}
