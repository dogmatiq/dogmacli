package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Call is a JSON-RPC method that has a response.
type Call struct {
	method

	Documentation       Documentation
	Params              Type
	RegistrationMethod  string
	RegistrationOptions Type
	Result              Type
	PartialResult       Type
	ErrorData           Type
}

func (b *builder) buildCall(in lowlevel.Request) *Call {
	return build(b, func(n *Call) {
		n.name = in.Method
		n.dir = methodDirection(in.Direction)

		n.Documentation = in.Documentation
		n.Params = b.buildType(in.Params)
		n.RegistrationMethod = in.RegistrationMethod
		n.RegistrationOptions = b.buildType(in.RegistrationOptions)
		n.Result = b.buildType(in.Result)
		n.PartialResult = b.buildType(in.PartialResult)
		n.ErrorData = b.buildType(in.ErrorData)
	})
}
