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

var (
	EOL = Token{typ: EndOfLine}
	EOF = Token{typ: EndOfFile}
)

type TokenType string

const (
	Unkonwn TokenType = "UNKONWN"

	// The operators
	EndOfLine      = "EOL"
	EndOfFile      = "EOF"
	OpLeftBracket  = "LBracket"
	OpRightBracket = "RBracket"
	OpLeftParen    = "LParen"
	OpRightParen   = "RParen"

	// The built-in reserved types
	TypeName   = "NAME"
	TypeString = "STRING"
	TypeInt    = "INT"

	// The built-in keywords
	KeyStr   = "str"
	KeyFn    = "fn"
	KeyPrint = "print"
)

type Token struct {
	raw string
	typ TokenType
}

// Create a new token instance and validate the token type
func NewToken(raw string) *Token {
	t := &Token{
		raw: raw,
		typ: classify(raw),
	}

	return t
}

// Show the raw token string
func (t *Token) String() string {
	return t.raw
}

// Show the token type
func (t *Token) Type() TokenType {
	return t.typ
}

func classify(raw string) TokenType {
	switch raw {
	case "(":
		return OpLeftParen
	case ")":
		return OpRightParen
	case "{":
		return OpLeftBracket
	case "}":
		return OpRightBracket
	case "str":
		return KeyStr
	case "fn":
		return KeyFn
	case "print":
		return KeyPrint
	default:
		switch {
		case reString.MatchString(raw):
			return TypeString
		case reInt.MatchString(raw):
			return TypeInt
		case reName.MatchString(raw):
			return TypeName
		}
	}

	return Unkonwn
}
