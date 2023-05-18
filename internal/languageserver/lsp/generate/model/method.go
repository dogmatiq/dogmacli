package model

// Method is a node that describes a JSON-RPC method.
type Method interface {
	Node

	MethodName() string
	Direction() MethodDirection
}

// MethodDirection indicates the direction that a JSON-RPC request is sent.
type MethodDirection int

type method struct {
	node
	name string
	dir  MethodDirection
}

func (m *method) MethodName() string {
	return m.name
}

func (m *method) Direction() MethodDirection {
	return m.dir
}

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
