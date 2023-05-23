package model

import "github.com/dogmatiq/dogmacli/internal/langserver/lsp/generate/model/internal/lowlevel"

// Call describes a JSON-RPC method that has a response.
type Call struct {
	methodDefNode

	ParamsType              Type
	RegistrationMethod      string
	RegistrationOptionsType Type

	ResultType        Type
	PartialResultType Type
	ErrorDataType     Type
}

func (b *builder) buildCall(in lowlevel.Request, out *Call) {
	out.name = in.Method
	out.docs = in.Documentation
	out.dir = methodDirection(in.Direction)

	out.ParamsType = b.buildType(in.Params)
	out.RegistrationMethod = in.RegistrationMethod
	out.RegistrationOptionsType = b.buildType(in.RegistrationOptions)
	out.ResultType = b.buildType(in.Result)
	out.PartialResultType = b.buildType(in.PartialResult)
	out.ErrorDataType = b.buildType(in.ErrorData)
}
