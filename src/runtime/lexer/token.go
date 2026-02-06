package lexer

// TokenType represents the type of a token.
type TokenType string

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	STRING TokenType = "STRING"

	// Operators
	ASSIGN  TokenType = "="
	DECLARE TokenType = ":="

	// Arithmetic operators
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"
	POWER    TokenType = "**"

	// Comparison operators
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	GT     TokenType = ">"
	LT_EQ  TokenType = "<="
	GT_EQ  TokenType = ">="

	// Delimiters
	COMMA  TokenType = ","
	LPAREN TokenType = "("
	RPAREN TokenType = ")"

	// Keywords
	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"
	NIL   TokenType = "NIL"
	MUT   TokenType = "MUT"
	AND   TokenType = "AND"
	OR    TokenType = "OR"
	NOT   TokenType = "NOT"
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"nil":   NIL,
	"mut":   MUT,
	"and":   AND,
	"or":    OR,
	"not":   NOT,
}

// LookupIdent returns the token type for an identifier.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
