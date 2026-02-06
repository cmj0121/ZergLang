package lexer

// Lexer performs lexical analysis on Zerg source code.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

// New creates a new Lexer for the given input.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	// TODO: Implement tokenization
	return Token{Type: EOF, Literal: "", Line: l.line, Column: l.column}
}

func (l *Lexer) readChar() {
	// TODO: Implement character reading
}
