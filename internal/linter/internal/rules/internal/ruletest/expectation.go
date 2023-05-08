package ruletest

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/dogmatiq/dogmacli/internal/linter/diagnostic"
	"golang.org/x/tools/go/packages"
)

// expectation describes an expected diagnostic.
type expectation struct {
	Severity diagnostic.Severity
	Message  string
	Pos      token.Position
}

// isLocatedWithin returns true if the given diagnostic is located at the
// expected position to satisfy this expectation.
func (e expectation) isLocatedWithin(d *diagnostic.Diagnostic) bool {
	if d.Begin.Filename != e.Pos.Filename {
		return false
	}

	if e.Pos.Line < d.Begin.Line {
		return false
	}

	if e.Pos.Line > d.End.Line {
		return false
	}

	return true
}

// isSatisfiedBy returns true if the given diagnostic satisfies the expectation.
func (e expectation) isSatisfiedBy(
	d *diagnostic.Diagnostic,
) bool {
	if d.Severity != e.Severity {
		return false
	}

	if d.Message != e.Message {
		return false
	}

	return e.isLocatedWithin(d)
}

// parseExpectations parses ruletest directives from the files in the given
// package.
func parseExpectations(pkg *packages.Package) ([]expectation, []error) {
	p := &expectationParser{
		Package: pkg,
	}

	for _, f := range pkg.Syntax {
		ast.Walk(p, f)
	}

	return p.Expectations, p.Errors
}

// expectationParser is an ast.Visitor that extracts ruletest directives from
// comments.
type expectationParser struct {
	Package      *packages.Package
	Expectations []expectation
	Errors       []error
}

func (p *expectationParser) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.File:
		for _, g := range n.Comments {
			for _, c := range g.List {
				p.parseDirective(c)
			}
		}
	}

	return nil
}

// expectationPattern is the regular expression used to parse ruletest
// directives from comments.
var expectationPattern = regexp.MustCompile(`^\/\/\s*ruletest:\s*\[([a-z]+)\]\s+(.+)$`)

// parseDirective parses a ruletest directive from a comment.
func (p *expectationParser) parseDirective(c *ast.Comment) {
	pos := p.Package.Fset.Position(c.Slash)

	matches := expectationPattern.FindStringSubmatch(c.Text)
	if len(matches) == 0 {
		if strings.Contains(c.Text, "ruletest") {
			p.Errors = append(
				p.Errors,
				fmt.Errorf(
					`invalid ruletest directive at %s, expected "ruletest: [<severity>] <message>"`,
					pos,
				),
			)
		}
		return
	}

	expect := expectation{
		Message: strings.TrimSpace(matches[2]),
		Pos:     pos,
	}

	sev, err := parseSeverity(matches[1])
	if err != nil {
		p.Errors = append(p.Errors, err)
		return
	}

	expect.Severity = sev
	p.Expectations = append(p.Expectations, expect)
}

// parseSeverity parses a diagnostic severity from a string.
func parseSeverity(s string) (diagnostic.Severity, error) {
	switch s {
	case "error":
		return diagnostic.Error, nil
	case "warning":
		return diagnostic.Warning, nil
	case "improvement":
		return diagnostic.Improvement, nil
	default:
		return 0, fmt.Errorf(`unrecognized severity %q, expected "error", "warning" or "improvement"`, s)
	}
}
