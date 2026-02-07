package lexer

// TokenType represents the type of a token.
type TokenType string

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENT        TokenType = "IDENT"
	INT          TokenType = "INT"
	FLOAT        TokenType = "FLOAT"
	STRING       TokenType = "STRING"
	INTERP_START TokenType = "INTERP_START" // Start of interpolated string "...{
	INTERP_MID   TokenType = "INTERP_MID"   // Middle part }...{
	INTERP_END   TokenType = "INTERP_END"   // End part }..."

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
	COMMA     TokenType = ","
	COLON     TokenType = ":"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	DOT       TokenType = "."
	ARROW     TokenType = "->"
	AMPERSAND TokenType = "&"

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

	// Builtins
	ASSERT TokenType = "ASSERT"

	// Low-level
	UNSAFE TokenType = "UNSAFE"
	ASM    TokenType = "ASM"

	// Module system
	IMPORT TokenType = "IMPORT"
	AS     TokenType = "AS"

	// Enum and Result types
	ENUM       TokenType = "ENUM"
	UNDERSCORE TokenType = "_"

	// Type checking
	IS TokenType = "IS"

	// Pipe for match alternatives
	PIPE TokenType = "|"

	// Range operators
	DOTDOT   TokenType = ".."
	DOTDOTEQ TokenType = "..="

	// Compound assignment operators
	PLUS_ASSIGN     TokenType = "+="
	MINUS_ASSIGN    TokenType = "-="
	ASTERISK_ASSIGN TokenType = "*="
	SLASH_ASSIGN    TokenType = "/="
	PERCENT_ASSIGN  TokenType = "%="
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
	"assert":   ASSERT,
	"unsafe":   UNSAFE,
	"asm":      ASM,
	"enum":     ENUM,
	"_":        UNDERSCORE,
	"is":       IS,
	"import":   IMPORT,
	"as":       AS,
}

// LookupIdent returns the token type for an identifier.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
