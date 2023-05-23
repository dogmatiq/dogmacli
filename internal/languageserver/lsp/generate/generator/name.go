package generator

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

// nameOf returns an identifier for the element described the given node.
func nameOf(n model.Node, suffix ...string) string {
	name, ok := tryNameOf(n, suffix...)
	if ok {
		return name
	}

	panic(fmt.Sprintf("%T has no Go identity", n))
}

func tryNameOf(n model.Node, suffix ...string) (string, bool) {
	// Any anonumous type that is directly referenced by an alias definition
	// inherits the name of that alias.
	if p, ok := n.Parent().(*model.Alias); ok {
		if p.UnderlyingType.IsAnonymous() {
			return tryNameOf(p)
		}
	}

	var v nameOfX
	model.VisitNode(n, &v)

	if v.N == "" {
		return "", false
	}

	return v.N + strings.Join(suffix, ""), true
}

type nameOfX struct {
	N string
}

func (v *nameOfX) VisitModel(n *model.Model)               {}
func (v *nameOfX) VisitCall(n *model.Call)                 { v.N = ident(n.Name()) }
func (v *nameOfX) VisitNotification(n *model.Notification) { v.N = ident(n.Name()) }
func (v *nameOfX) VisitBool(n *model.Bool)                 { v.N = "Bool" }
func (v *nameOfX) VisitDecimal(n *model.Decimal)           { v.N = "Decimal" }
func (v *nameOfX) VisitString(n *model.String)             { v.N = "String" }
func (v *nameOfX) VisitInteger(n *model.Integer)           { v.N = "Int" }
func (v *nameOfX) VisitUInteger(n *model.UInteger)         { v.N = "UInt" }
func (v *nameOfX) VisitDocumentURI(n *model.DocumentURI)   { v.N = "DocumentURI" }
func (v *nameOfX) VisitURI(n *model.URI)                   { v.N = "URI" }
func (v *nameOfX) VisitNull(n *model.Null)                 {}
func (v *nameOfX) VisitReference(n *model.Reference)       { v.N = nameOf(n.Target) }
func (v *nameOfX) VisitArray(n *model.Array)               { v.N = nameOf(n.ElementType, "Array") }
func (v *nameOfX) VisitMap(n *model.Map)                   { v.N = nameOf(n.ValueType, "Map") }
func (v *nameOfX) VisitAnd(n *model.And)                   { v.N = scopeOf(n) }
func (v *nameOfX) VisitOr(n *model.Or)                     { v.N = scopeOf(n) }
func (v *nameOfX) VisitTuple(n *model.Tuple)               { v.N = scopeOf(n) }
func (v *nameOfX) VisitStructLit(n *model.StructLit)       { v.N = scopeOf(n) }
func (v *nameOfX) VisitStringLit(n *model.StringLit)       {}
func (v *nameOfX) VisitAlias(n *model.Alias)               { v.N = ident(n.Name()) }
func (v *nameOfX) VisitEnum(n *model.Enum)                 { v.N = ident(n.Name()) }
func (v *nameOfX) VisitEnumMember(n *model.EnumMember)     { v.N = ident(n.Name) + nameOf(n.Parent()) }
func (v *nameOfX) VisitStruct(n *model.Struct)             { v.N = ident(n.Name()) }
func (v *nameOfX) VisitProperty(n *model.Property)         { v.N = ident(n.Name) }

// scopeOf returns an identifier for the "scopeOf" that n is within. It is used
// to give anonymous types a name based on where they are defined.
func scopeOf(n model.Node) string {
	v := scopeOfX{
		Child: n,
	}

	for v.Child.HasParent() {
		p := v.Child.Parent()
		model.VisitNode(p, &v)
		v.Child = p
	}

	return v.Name
}

type scopeOfX struct {
	Child model.Node
	Name  string
}

func (v *scopeOfX) VisitCall(n *model.Call) {
	switch v.Child {
	case n.ParamsType:
		v.push("Params")
	case n.RegistrationOptionsType:
		v.push("RegistrationOptions")
	case n.ResultType:
		v.push("Result")
	case n.PartialResultType:
		v.push("PartialResult")
	case n.ErrorDataType:
		v.push("Error")
	default:
		panic("child not found in parent")
	}

	v.push(ident(n.Name()))
}

func (v *scopeOfX) VisitNotification(n *model.Notification) {
	switch v.Child {
	case n.ParamsType:
		v.push("Params")
	case n.RegistrationOptionsType:
		v.push("RegistrationOptions")
	default:
		panic("child not found in parent")
	}

	v.push(ident(n.Name()))
}

func (v *scopeOfX) VisitOr(n *model.Or) {
	letter := rune('A')
	for _, t := range n.Types {
		if t == v.Child {
			v.push(string(letter))
			return
		}
		if t.IsAnonymous() {
			letter++
		}
	}
	panic("child not found in parent")
}

func (v *scopeOfX) VisitTuple(n *model.Tuple) {
	for i, t := range n.Types {
		if t == v.Child {
			v.push(ordinals[i])
			return
		}
	}
	panic("child not found in parent")
}

func (v *scopeOfX) VisitAlias(n *model.Alias) {
	v.push(ident(n.Name()))
}

func (v *scopeOfX) VisitStruct(n *model.Struct) {
	v.push(ident(n.Name()))
}

func (v *scopeOfX) VisitProperty(n *model.Property) {
	v.push(ident(n.Name))
}

// These node types may have anonymous types declared within them.
func (v *scopeOfX) VisitModel(n *model.Model)         {}
func (v *scopeOfX) VisitReference(n *model.Reference) {}
func (v *scopeOfX) VisitArray(n *model.Array)         {}
func (v *scopeOfX) VisitMap(n *model.Map)             {}
func (v *scopeOfX) VisitStructLit(n *model.StructLit) {}

// These node types are not expected to have anonymous types declared within them.
func (v *scopeOfX) VisitBool(n *model.Bool)               { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitDecimal(n *model.Decimal)         { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitString(n *model.String)           { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitInteger(n *model.Integer)         { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitUInteger(n *model.UInteger)       { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitDocumentURI(n *model.DocumentURI) { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitURI(n *model.URI)                 { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitNull(n *model.Null)               { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitStringLit(n *model.StringLit)     { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitAnd(n *model.And)                 { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitEnum(n *model.Enum)               { panic("unexpected anonymous type") }
func (v *scopeOfX) VisitEnumMember(n *model.EnumMember)   { panic("unexpected anonymous type") }

func (v *scopeOfX) push(name string) {
	v.Name = ident(name) + v.Name
}

var ordinals = [...]string{
	"First",
	"Second",
	"Third",
	"Fourth",
	"Fifth",
	"Sixth",
	"Seventh",
	"Eighth",
	"Ninth",
	"Tenth",
}

// ident returns a ident ident containing the given parts.
func ident(name string) string {
	if name[0] == '_' {
		name = name[1:]
		return unexported(name)
	}
	return exported(name)
}

// exported returns a normalized exported identifier containing the given parts.
func exported(name string) string {
	var id string

	for _, x := range strings.Split(name, "/") {
		if x != "$" {
			id += strings.Title(x)
		}
	}

	if x, ok := strings.CutSuffix(id, "Id"); ok {
		id = x + "ID"
	}
	if x, ok := strings.CutSuffix(id, "Ids"); ok {
		id = x + "IDs"
	}
	if x, ok := strings.CutSuffix(id, "Uri"); ok {
		id = x + "URI"
	}

	return id
}

// unexported returns a normalized unexported identifier containing the given
// parts.
func unexported(name string) string {
	id := exported(name)
	return strings.ToLower(id[:1]) + id[1:]
}
