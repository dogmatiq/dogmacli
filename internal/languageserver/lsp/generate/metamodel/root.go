package metamodel

import (
	_ "embed"
	"encoding/json"

	"golang.org/x/exp/slices"
)

// Root is the root of the model.
type Root struct {
	ServerAPI  []RPC
	ClientAPI  []RPC
	NamedTypes []NamedType
}

//go:embed metamodel-3.17.0.json
var data []byte

// rootJSON is the JSON representation of the meta-model.
type rootJSON struct {
	Requests      []requestJSON      `json:"requests"`
	Notifications []notificationJSON `json:"notifications"`
	Structures    []structureJSON    `json:"structures"`
	Enumerations  []enumerationJSON  `json:"enumerations"`
	TypeAliases   []typeAliasJSON    `json:"typeAliases"`
}

// Get returns the root node of the meta-model.
func Get() Root {
	var j rootJSON
	if err := json.Unmarshal(data, &j); err != nil {
		panic(err)
	}

	// Pre-construct pointers to the named types so that can be obtained
	// wherever they are referenced.
	var (
		root         Root
		named        = map[string]NamedType{}
		structures   = map[string]*Structure{}
		enumerations = map[string]*Enumeration{}
		typeAliases  = map[string]*TypeAlias{}
	)

	for _, jt := range j.Structures {
		t := &Structure{}
		named[jt.Name] = t
		structures[jt.Name] = t
		root.NamedTypes = append(root.NamedTypes, t)
	}

	for _, jt := range j.Enumerations {
		t := &Enumeration{}
		named[jt.Name] = t
		enumerations[t.Name] = t
		root.NamedTypes = append(root.NamedTypes, t)
	}

	for _, jt := range j.TypeAliases {
		t := &TypeAlias{}
		named[jt.Name] = t
		typeAliases[t.Name] = t
		root.NamedTypes = append(root.NamedTypes, t)
	}

	for _, t := range j.Structures {
		populateStructure(named, t)
	}

	for _, t := range j.Enumerations {
		populateEnumeration(named, t)
	}

	for _, t := range j.TypeAliases {
		populateTypeAlias(named, t)
	}

	for _, rpc := range j.Requests {
		if rpc.Direction == "clientToServer" {
			root.ClientAPI = append(root.ClientAPI, newRequest(named, rpc))
		} else {
			root.ServerAPI = append(root.ServerAPI, newRequest(named, rpc))
		}
	}

	for _, rpc := range j.Notifications {
		if rpc.Direction == "clientToServer" {
			root.ClientAPI = append(root.ClientAPI, newNotification(named, rpc))
		} else {
			root.ServerAPI = append(root.ServerAPI, newNotification(named, rpc))
		}
	}

	slices.SortFunc(
		root.ServerAPI,
		func(a, b RPC) bool {
			return a.MethodName() < b.MethodName()
		},
	)

	slices.SortFunc(
		root.ClientAPI,
		func(a, b RPC) bool {
			return a.MethodName() < b.MethodName()
		},
	)

	slices.SortFunc(
		root.NamedTypes,
		func(a, b NamedType) bool {
			return a.TypeName() < b.TypeName()
		},
	)

	return root
}
