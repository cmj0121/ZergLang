package token

var (
	EndOfLine = Token{typ: EOL}
	EndOfFile = Token{typ: EOF}
)

type Token struct {
	raw string
	typ Type
}

// Create a new token instance and validate the token type
func NewToken(raw string) *Token {
	t := &Token{
		raw: raw,
		typ: NewType(raw),
	}

	return t
}

// Show the raw token string
func (t *Token) String() string {
	return t.raw
}

// Show the token type
func (t *Token) Type() Type {
	return t.typ
}
