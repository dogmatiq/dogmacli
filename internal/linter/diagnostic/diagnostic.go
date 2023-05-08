package diagnostic

import (
	"go/token"
)

// Diagnostic represents a diagnostic message.
type Diagnostic struct {
	Severity   Severity
	Begin, End token.Position
	Message    string
	Suggestion *Change
}

// SuggestChange adds a set of edits that the user may perform to resolve this
// diagnostic.
func (d *Diagnostic) SuggestChange(message string, edits ...Edit) *Diagnostic {
	d.Suggestion = &Change{message, edits}
	return d
}

// Severity is an enumeration describing the severity of a diagnostic.
type Severity int

const (
	// Error indicates that the diagnostic describes usage of Dogma that will
	// not function correctly.
	Error Severity = iota

	// Warning indicates that the diagnostic describes usage of Dogma that is
	// not recommended.
	Warning

	// Improvement indicates that the describes usage of Dogma that may be
	// improved, but is not problematic.
	Improvement
)

func (s Severity) String() string {
	switch s {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Improvement:
		return "improvement"
	default:
		return "unknown"
	}
}

// DapperString returns the string representation of the diagnostic severity for
// use with dogmatiq/dapper.
func (s Severity) DapperString() string {
	return s.String()
}

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
