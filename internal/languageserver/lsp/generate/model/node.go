package model

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model/internal/lowlevel"

// Node is an interface for all nodes within the meta-model.
type Node interface {
	// Parent returns the node's parent, if any.
	Parent() (Node, bool)
	setParent(Node)
	acceptVisitor(Visitor)
}

type node struct {
	par Node
}

func (n node) Parent() (Node, bool) {
	return n.par, n.par != nil
}

func (n *node) setParent(p Node) {
	n.par = p
}

// Documentation is a container for documentation about a node.
type Documentation = lowlevel.Documentation

// Visitor is an interface for node-type specific logic.
type Visitor interface {
	VisitModel(*Model)
	VisitCall(*Call)
	VisitNotification(*Notification)
	VisitBool(*Bool)
	VisitDecimal(*Decimal)
	VisitString(*String)
	VisitInteger(*Integer)
	VisitUInteger(*UInteger)
	VisitDocumentURI(*DocumentURI)
	VisitURI(*URI)
	VisitNull(*Null)
	VisitReference(*Reference)
	VisitArray(*Array)
	VisitMap(*Map)
	VisitAnd(*And)
	VisitOr(*Or)
	VisitTuple(*Tuple)
	VisitStructLit(*StructLit)
	VisitStringLit(*StringLit)
	VisitAlias(*Alias)
	VisitEnum(*Enum)
	VisitEnumMember(*EnumMember)
	VisitStruct(*Struct)
	VisitProperty(*Property)
}

// Visit dispatches to the appropriate method on the given visitor base on the
// concrete type of n.
func Visit(n Node, v Visitor) {
	n.acceptVisitor(v)
}

func (n *Model) acceptVisitor(v Visitor)        { v.VisitModel(n) }
func (n *Call) acceptVisitor(v Visitor)         { v.VisitCall(n) }
func (n *Notification) acceptVisitor(v Visitor) { v.VisitNotification(n) }
func (n *Bool) acceptVisitor(v Visitor)         { v.VisitBool(n) }
func (n *Decimal) acceptVisitor(v Visitor)      { v.VisitDecimal(n) }
func (n *String) acceptVisitor(v Visitor)       { v.VisitString(n) }
func (n *Integer) acceptVisitor(v Visitor)      { v.VisitInteger(n) }
func (n *UInteger) acceptVisitor(v Visitor)     { v.VisitUInteger(n) }
func (n *DocumentURI) acceptVisitor(v Visitor)  { v.VisitDocumentURI(n) }
func (n *URI) acceptVisitor(v Visitor)          { v.VisitURI(n) }
func (n *Null) acceptVisitor(v Visitor)         { v.VisitNull(n) }
func (n *Reference) acceptVisitor(v Visitor)    { v.VisitReference(n) }
func (n *Array) acceptVisitor(v Visitor)        { v.VisitArray(n) }
func (n *Map) acceptVisitor(v Visitor)          { v.VisitMap(n) }
func (n *And) acceptVisitor(v Visitor)          { v.VisitAnd(n) }
func (n *Or) acceptVisitor(v Visitor)           { v.VisitOr(n) }
func (n *Tuple) acceptVisitor(v Visitor)        { v.VisitTuple(n) }
func (n *StructLit) acceptVisitor(v Visitor)    { v.VisitStructLit(n) }
func (n *StringLit) acceptVisitor(v Visitor)    { v.VisitStringLit(n) }
func (n *Alias) acceptVisitor(v Visitor)        { v.VisitAlias(n) }
func (n *Enum) acceptVisitor(v Visitor)         { v.VisitEnum(n) }
func (n *EnumMember) acceptVisitor(v Visitor)   { v.VisitEnumMember(n) }
func (n *Struct) acceptVisitor(v Visitor)       { v.VisitStruct(n) }
func (n *Property) acceptVisitor(v Visitor)     { v.VisitProperty(n) }
