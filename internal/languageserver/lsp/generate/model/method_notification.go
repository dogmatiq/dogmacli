package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Notification is a JSON-RPC method that has no response.
type Notification struct {
	method

	Documentation       Documentation
	Params              Type
	RegistrationMethod  string
	RegistrationOptions Type
}

func (b *builder) buildNotification(in lowlevel.Notification) *Notification {
	return build(b, func(n *Notification) {
		n.name = in.Method
		n.dir = methodDirection(in.Direction)

		n.Documentation = in.Documentation
		n.Params = b.buildType(in.Params)
		n.RegistrationMethod = in.RegistrationMethod
		n.RegistrationOptions = b.buildType(in.RegistrationOptions)
	})
}
