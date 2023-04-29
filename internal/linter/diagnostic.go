package linter

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Reporter allows a rule to report diagnostics.
type Reporter struct {
	diagnostics []*Diagnostic
}

// Report reports a diagnostic about the given AST node.
func (r *Reporter) Report(
	s DiagnosticSeverity,
	n ast.Node,
	format string, args ...any,
) *Diagnostic {
	d := &Diagnostic{
		Begin:   n.Pos(),
		End:     n.End(),
		Message: fmt.Sprintf(format, args...),
	}

	r.diagnostics = append(r.diagnostics, d)

	return d
}

// Error reports an error about the given AST node.
func (r *Reporter) Error(
	n ast.Node,
	format string, args ...any,
) *Diagnostic {
	return r.Report(Error, n, format, args...)
}

// Warning reports a warning about the given AST node.
func (r *Reporter) Warning(
	n ast.Node,
	format string, args ...any,
) *Diagnostic {
	return r.Report(Warning, n, format, args...)
}

// Improvement reports an improvement that can be made to the given AST node.
func (r *Reporter) Improvement(
	n ast.Node,
	format string, args ...any,
) *Diagnostic {
	return r.Report(Improvement, n, format, args...)
}

// Diagnostic represents a diagnostic message.
type Diagnostic struct {
	Severity   DiagnosticSeverity
	Begin, End token.Pos
	Message    string
	Suggestion *Change
}

// SuggestChange adds a set of edits that the user may perform to resolve this
// diagnostic.
func (d *Diagnostic) SuggestChange(message string, edits ...Edit) *Diagnostic {
	d.Suggestion = &Change{message, edits}
	return d
}

// DiagnosticSeverity is an enumeration of descibing the severity of a
// diagnostic.
type DiagnosticSeverity int

const (
	// Error indicates that the diagnostic describes usage of Dogma that will
	// not function correctly.
	Error DiagnosticSeverity = iota

	// Warning indicates that the diagnostic describes usage of Dogma that is
	// not recommended.
	Warning

	// Improvement indicates that the describes usage of Dogma that may be
	// improved, but is not problematic.
	Improvement
)

// Change is related set of edits to a document.
type Change struct {
	Message string
	Edits   []Edit
}

// Edit is a single contiguous edit within a document.
type Edit struct {
	Begin, End token.Pos
	Text       string
}
