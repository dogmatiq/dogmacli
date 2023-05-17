package lsp

type Optional[T any] struct {
	value *T
}

func (o Optional[T]) Get() (T, bool) {
	if o.value != nil {
		return *o.value, true
	}

	var zero T
	return zero, false
}
