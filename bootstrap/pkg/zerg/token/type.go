// The valid token types.
package token

import (
	"regexp"
)

var (
	// The build-in type regexp pattern
	reName   = regexp.MustCompile(`^[a-zA-Z_]\w*$`)
	reString = regexp.MustCompile(`^".*"$`)
	reInt    = regexp.MustCompile(`^\d+$`)
)

//go:generate stringer -type=Type
type Type int

const (
	Unknown Type = iota

	EOF
	EOL

	// The operators
	LBracket
	RBracket
	LParen
	RParen
	LBrace
	RBrace

	Add
	Sub
	Mul
	Div
	Mod

	// The composite operators
	Arrow

	// The build-in type and reserved keywords
	Name
	String
	Int
	Fn
	Str
	Print
)

// Identify the raw token string and return the token type.
func NewType(raw string) Type {
	mapping := map[string]Type{
		"(":     LParen,
		")":     RParen,
		"{":     LBrace,
		"}":     RBrace,
		"[":     LBracket,
		"]":     RBracket,
		"+":     Add,
		"-":     Sub,
		"*":     Mul,
		"/":     Div,
		"%":     Mod,
		"->":    Arrow,
		"str":   Str,
		"fn":    Fn,
		"print": Print,
	}

	switch typ, ok := mapping[raw]; {
	case ok:
		return typ
	default:
		return ComplexType(raw)
	}
}

// Identify the complex token type.
func ComplexType(raw string) Type {
	switch {
	case reString.MatchString(raw):
		return String
	case reInt.MatchString(raw):
		return Int
	case reName.MatchString(raw):
		return Name
	default:
		return Unknown
	}
}
