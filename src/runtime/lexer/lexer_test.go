package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `x := 42
y := "hello"
true false nil
()`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "x"},
		{DECLARE, ":="},
		{INT, "42"},
		{IDENT, "y"},
		{DECLARE, ":="},
		{STRING, "hello"},
		{TRUE, "true"},
		{FALSE, "false"},
		{NIL, "nil"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIntegerWithUnderscores(t *testing.T) {
	input := `1_000_000`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != INT {
		t.Fatalf("expected INT, got %s", tok.Type)
	}

	if tok.Literal != "1_000_000" {
		t.Fatalf("expected 1_000_000, got %s", tok.Literal)
	}
}

func TestComments(t *testing.T) {
	input := `# this is a comment
x := 42  # inline comment
y`

	l := New(input)

	tok := l.NextToken()
	if tok.Type != IDENT || tok.Literal != "x" {
		t.Fatalf("expected x, got %s", tok.Literal)
	}

	l.NextToken() // :=
	l.NextToken() // 42

	tok = l.NextToken()
	if tok.Type != IDENT || tok.Literal != "y" {
		t.Fatalf("expected y, got %s", tok.Literal)
	}
}

func TestStringEscapes(t *testing.T) {
	input := `"hello\nworld"`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != STRING {
		t.Fatalf("expected STRING, got %s", tok.Type)
	}

	// Escape sequences are now converted to actual characters
	expected := "hello\nworld"
	if tok.Literal != expected {
		t.Fatalf("expected %q, got %q", expected, tok.Literal)
	}
}

func TestStringEscapesAll(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello\nworld"`, "hello\nworld"},
		{`"hello\tworld"`, "hello\tworld"},
		{`"hello\\world"`, "hello\\world"},
		{`"hello\"world"`, "hello\"world"},
		{`"hello\rworld"`, "hello\rworld"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != STRING {
			t.Errorf("for %s: expected STRING, got %s", tt.input, tok.Type)
			continue
		}

		if tok.Literal != tt.expected {
			t.Errorf("for %s: expected %q, got %q", tt.input, tt.expected, tok.Literal)
		}
	}
}

func TestAssignmentTokens(t *testing.T) {
	input := `mut x := 10
x = 20
a, b = b, a`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{MUT, "mut"},
		{IDENT, "x"},
		{DECLARE, ":="},
		{INT, "10"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT, "20"},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{ASSIGN, "="},
		{IDENT, "b"},
		{COMMA, ","},
		{IDENT, "a"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
