package metamodel

// RPC is an interface for JSON-RPC calls and notifications.
type RPC interface {
	MethodName() string
}

// Request defines a JSON-RPC call (request/response).
type Request struct {
	Method              string
	Documentation       string
	Parameters          Type
	Result              Type
	PartialResult       Type
	RegistrationMethod  string
	RegistrationOptions Type
}

func newRequest(named map[string]NamedType, rpc requestJSON) Request {
	return Request{
		Method:              rpc.Method,
		Documentation:       rpc.Documentation,
		Parameters:          newType(named, rpc.Params),
		Result:              newType(named, rpc.Result),
		PartialResult:       newType(named, rpc.PartialResult),
		RegistrationMethod:  rpc.RegistrationMethod,
		RegistrationOptions: newType(named, rpc.RegistrationOptions),
	}
}

// MethodName returns the name of the RPC method.
func (c Request) MethodName() string {
	return c.Method
}

// Notification defines a JSON-RPC notification.
type Notification struct {
	Method              string
	Documentation       string
	Parameters          Type
	RegistrationMethod  string
	RegistrationOptions Type
}

func newNotification(named map[string]NamedType, rpc notificationJSON) Notification {
	return Notification{
		Method:              rpc.Method,
		Documentation:       rpc.Documentation,
		Parameters:          newType(named, rpc.Params),
		RegistrationMethod:  rpc.RegistrationMethod,
		RegistrationOptions: newType(named, rpc.RegistrationOptions),
	}
}

// MethodName returns the name of the RPC method.
func (n Notification) MethodName() string {
	return n.Method
}

// Request defines a JSON-RPC call (requestJSON/response).
type requestJSON struct {
	Method              string    `json:"method"`
	Documentation       string    `json:"documentation"`
	Direction           string    `json:"messageDirection"`
	Params              *typeJSON `json:"params"`
	Result              *typeJSON `json:"result"`
	PartialResult       *typeJSON `json:"partialResult"`
	RegistrationMethod  string    `json:"registrationMethod"`
	RegistrationOptions *typeJSON `json:"registrationOptions"`
}

// notificationJSON defines a JSON-RPC notificationJSON.
type notificationJSON struct {
	Method              string    `json:"method"`
	Documentation       string    `json:"documentation"`
	Direction           string    `json:"messageDirection"`
	Params              *typeJSON `json:"params"`
	RegistrationMethod  string    `json:"registrationMethod"`
	RegistrationOptions *typeJSON `json:"registrationOptions"`
}
