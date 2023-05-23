package generator

import (
	"reflect"

	"github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"
)

// kindOf returns the kind of type produced by n.
func kindOf(n model.Node) reflect.Kind {
	var v kindOfX
	model.VisitNode(n, &v)
	return v.K
}

type kindOfX struct {
	K reflect.Kind
}

func (v *kindOfX) VisitModel(n *model.Model)               {}
func (v *kindOfX) VisitCall(n *model.Call)                 {}
func (v *kindOfX) VisitNotification(n *model.Notification) {}
func (v *kindOfX) VisitBool(n *model.Bool)                 { v.K = reflect.Bool }
func (v *kindOfX) VisitDecimal(n *model.Decimal)           { v.K = reflect.Float64 }
func (v *kindOfX) VisitString(n *model.String)             { v.K = reflect.String }
func (v *kindOfX) VisitInteger(n *model.Integer)           { v.K = reflect.Int32 }
func (v *kindOfX) VisitUInteger(n *model.UInteger)         { v.K = reflect.Uint32 }
func (v *kindOfX) VisitDocumentURI(n *model.DocumentURI)   { v.K = reflect.Struct }
func (v *kindOfX) VisitURI(n *model.URI)                   { v.K = reflect.Struct }
func (v *kindOfX) VisitNull(n *model.Null)                 {}
func (v *kindOfX) VisitReference(n *model.Reference)       { v.K = kindOf(n.Target) }
func (v *kindOfX) VisitArray(n *model.Array)               { v.K = reflect.Slice }
func (v *kindOfX) VisitMap(n *model.Map)                   { v.K = reflect.Map }
func (v *kindOfX) VisitAnd(n *model.And)                   { v.K = reflect.Struct }
func (v *kindOfX) VisitOr(n *model.Or)                     { v.K = reflect.Interface }
func (v *kindOfX) VisitTuple(n *model.Tuple)               { v.K = reflect.Array }
func (v *kindOfX) VisitStructLit(n *model.StructLit)       { v.K = reflect.Struct }
func (v *kindOfX) VisitStringLit(n *model.StringLit)       {}
func (v *kindOfX) VisitAlias(n *model.Alias)               { v.K = kindOf(n.UnderlyingType) }
func (v *kindOfX) VisitEnum(n *model.Enum)                 { v.K = kindOf(n.UnderlyingType) }
func (v *kindOfX) VisitEnumMember(n *model.EnumMember)     { v.K = kindOf(n.Parent()) }
func (v *kindOfX) VisitStruct(n *model.Struct)             { v.K = reflect.Struct }
func (v *kindOfX) VisitProperty(n *model.Property)         { v.K = kindOf(n.Type) }
