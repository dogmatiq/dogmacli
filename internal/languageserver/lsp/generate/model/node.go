package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Node is an interface for all nodes within the meta-model.
type Node interface {
	HasParent() bool
	Parent() Node
	setParent(Node)

	acceptVisitor(NodeVisitor)
}

type node struct {
	parent Node
}

func (n *node) HasParent() bool {
	return n.parent != nil
}

func (n *node) Parent() Node {
	if n.parent == nil {
		panic("root has no parent")
	}
	return n.parent
}

func (n *node) setParent(p Node) {
	n.parent = p
}

// Documentation is a container for documentation about a node.
type Documentation = lowlevel.Documentation

// NodeVisitor is an interface for node-type specific logic.
type NodeVisitor interface {
	TypeVisitor

	VisitModel(*Model)
	VisitCall(*Call)
	VisitNotification(*Notification)
	VisitAlias(*Alias)
	VisitEnum(*Enum)
	VisitEnumMember(*EnumMember)
	VisitStruct(*Struct)
	VisitProperty(*Property)
}

// VisitNode dispatches to the appropriate method on the given visitor base on
// the concrete type of n.
func VisitNode(n Node, v NodeVisitor) {
	n.acceptVisitor(v)
}

func (n *Model) acceptVisitor(v NodeVisitor)        { v.VisitModel(n) }
func (n *Call) acceptVisitor(v NodeVisitor)         { v.VisitCall(n) }
func (n *Notification) acceptVisitor(v NodeVisitor) { v.VisitNotification(n) }
func (n *Bool) acceptVisitor(v NodeVisitor)         { v.VisitBool(n) }
func (n *Decimal) acceptVisitor(v NodeVisitor)      { v.VisitDecimal(n) }
func (n *String) acceptVisitor(v NodeVisitor)       { v.VisitString(n) }
func (n *Integer) acceptVisitor(v NodeVisitor)      { v.VisitInteger(n) }
func (n *UInteger) acceptVisitor(v NodeVisitor)     { v.VisitUInteger(n) }
func (n *DocumentURI) acceptVisitor(v NodeVisitor)  { v.VisitDocumentURI(n) }
func (n *URI) acceptVisitor(v NodeVisitor)          { v.VisitURI(n) }
func (n *Null) acceptVisitor(v NodeVisitor)         { v.VisitNull(n) }
func (n *Reference) acceptVisitor(v NodeVisitor)    { v.VisitReference(n) }
func (n *Array) acceptVisitor(v NodeVisitor)        { v.VisitArray(n) }
func (n *Map) acceptVisitor(v NodeVisitor)          { v.VisitMap(n) }
func (n *And) acceptVisitor(v NodeVisitor)          { v.VisitAnd(n) }
func (n *Or) acceptVisitor(v NodeVisitor)           { v.VisitOr(n) }
func (n *Tuple) acceptVisitor(v NodeVisitor)        { v.VisitTuple(n) }
func (n *StructLit) acceptVisitor(v NodeVisitor)    { v.VisitStructLit(n) }
func (n *StringLit) acceptVisitor(v NodeVisitor)    { v.VisitStringLit(n) }
func (n *Alias) acceptVisitor(v NodeVisitor)        { v.VisitAlias(n) }
func (n *Enum) acceptVisitor(v NodeVisitor)         { v.VisitEnum(n) }
func (n *EnumMember) acceptVisitor(v NodeVisitor)   { v.VisitEnumMember(n) }
func (n *Struct) acceptVisitor(v NodeVisitor)       { v.VisitStruct(n) }
func (n *Property) acceptVisitor(v NodeVisitor)     { v.VisitProperty(n) }
