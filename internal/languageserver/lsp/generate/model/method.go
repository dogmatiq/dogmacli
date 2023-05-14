package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Method is an interface for JSON-RPC methods.
type Method interface {
	Name() string
	accept(MethodVisitor)
}

type (
	// MethodCommon contains the common fields of Call and Notification.
	MethodCommon struct {
		MethodName          string
		Documentation       Documentation
		Params              Type
		RegistrationMethod  string
		RegistrationOptions Type
	}

	// Call is a JSON-RPC method that has a response.
	Call struct {
		MethodCommon

		Result        Type
		PartialResult Type
		ErrorData     Type
	}

	// Notification is a JSON-RPC method that has no response.
	Notification struct {
		MethodCommon
	}
)

// Name returns the method name.
func (m Call) Name() string {
	return m.MethodName
}

// Name returns the method name.
func (m Notification) Name() string {
	return m.MethodName
}

func (b *builder) call(in lowlevel.Request) Call {
	return Call{
		MethodCommon: MethodCommon{
			MethodName:          in.Method,
			Documentation:       in.Documentation,
			Params:              b.typeRef(in.Params),
			RegistrationMethod:  in.RegistrationMethod,
			RegistrationOptions: b.typeRef(in.RegistrationOptions),
		},
		Result:        b.typeRef(in.Result),
		PartialResult: b.typeRef(in.PartialResult),
		ErrorData:     b.typeRef(in.ErrorData),
	}
}

func (b *builder) notification(in lowlevel.Notification) Notification {
	return Notification{
		MethodCommon: MethodCommon{
			MethodName:          in.Method,
			Documentation:       in.Documentation,
			Params:              b.typeRef(in.Params),
			RegistrationMethod:  in.RegistrationMethod,
			RegistrationOptions: b.typeRef(in.RegistrationOptions),
		},
	}
}

// MethodVisitor provides logic specific to each Method implementation.
type MethodVisitor interface {
	Call(Call)
	Notification(Notification)
}

// VisitMethod dispatches to v based on the concrete type of m.
func VisitMethod(m Method, v MethodVisitor) {
	m.accept(v)
}

// MethodTransform produces a value of type T from a TypeDef.
type MethodTransform[T any] interface {
	Call(Call) T
	Notification(Notification) T
}

// ApplyMethodTransform transforms m to a value of type T using x.
func ApplyMethodTransform[T any](
	m Method,
	x MethodTransform[T],
) T {
	v := &methodX[T]{X: x}
	VisitMethod(m, v)
	return v.V
}

type methodX[T any] struct {
	X MethodTransform[T]
	V T
}

func (m Call) accept(v MethodVisitor)         { v.Call(m) }
func (m Notification) accept(v MethodVisitor) { v.Notification(m) }

func (v *methodX[T]) Call(m Call)                 { v.V = v.X.Call(m) }
func (v *methodX[T]) Notification(m Notification) { v.V = v.X.Notification(m) }
