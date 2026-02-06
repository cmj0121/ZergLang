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
	COMMA    TokenType = ","
	COLON    TokenType = ":"
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	DOT      TokenType = "."
	ARROW    TokenType = "->"

	// Keywords
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	NIL      TokenType = "NIL"
	MUT      TokenType = "MUT"
	AND      TokenType = "AND"
	OR       TokenType = "OR"
	NOT      TokenType = "NOT"
	FN       TokenType = "FN"
	RETURN   TokenType = "RETURN"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	FOR      TokenType = "FOR"
	IN       TokenType = "IN"
	BREAK    TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"
	NOP      TokenType = "NOP"

	// Match (for future)
	MATCH     TokenType = "MATCH"
	FAT_ARROW TokenType = "=>"

	// Class-related keywords
	CLASS  TokenType = "CLASS"
	IMPL   TokenType = "IMPL"
	THIS   TokenType = "THIS"
	PUB    TokenType = "PUB"
	STATIC TokenType = "STATIC"

	// Spec-related keywords
	SPEC TokenType = "SPEC"
	SELF TokenType = "SELF"
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
	"mut":      MUT,
	"and":      AND,
	"or":       OR,
	"not":      NOT,
	"fn":       FN,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"in":       IN,
	"break":    BREAK,
	"continue": CONTINUE,
	"nop":      NOP,
	"match":    MATCH,
	"class":    CLASS,
	"impl":     IMPL,
	"this":     THIS,
	"pub":      PUB,
	"static":   STATIC,
	"spec":     SPEC,
	"Self":     SELF,
}

// LookupIdent returns the token type for an identifier.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
