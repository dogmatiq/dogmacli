package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Notification describes a JSON-RPC method that has no response.
type Notification struct {
	methodDefNode

	ParamsType              Type
	RegistrationMethod      string
	RegistrationOptionsType Type
}

func (b *builder) buildNotification(in lowlevel.Notification, out *Notification) {
	out.name = in.Method
	out.docs = in.Documentation
	out.dir = methodDirection(in.Direction)

	out.ParamsType = b.buildType(in.Params)
	out.RegistrationMethod = in.RegistrationMethod
	out.RegistrationOptionsType = b.buildType(in.RegistrationOptions)
}
