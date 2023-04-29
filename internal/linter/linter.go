package linter

// Linter reports problems with Dogma applications and handlers.
type Linter struct {
	Checks []func(*Context)
}

func (l *Linter) Lint() {
}
