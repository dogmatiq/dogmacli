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
