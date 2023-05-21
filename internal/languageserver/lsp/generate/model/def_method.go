package model

// MethodDef describes a JSON-RPC method.
type MethodDef interface {
	Def

	Direction() MethodDirection
}

// methodDefNode provides implementation common to all types that implement
// MethodDef.
type methodDefNode struct {
	defNode

	dir MethodDirection
}

func (m *methodDefNode) Direction() MethodDirection {
	return m.dir
}

// MethodDirection indicates the direction that a JSON-RPC request is sent.
type MethodDirection int

const (
	// HandledByLanguageServer indicates that a JSON-RPC method request is sent
	// from the IDE to the language server.
	HandledByLanguageServer MethodDirection = iota

	// HandledByIDE indicates that a JSON-RPC method request is sent from the
	// language server to the IDE.
	HandledByIDE
)

// methodDirection converts a low-level message direction to a MethodDirection.
func methodDirection(dir string) MethodDirection {
	if dir == "clientToServer" {
		return HandledByLanguageServer
	}
	return HandledByIDE
}
