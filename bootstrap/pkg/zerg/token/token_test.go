package token

import (
	"testing"
)

func TestToken(t *testing.T) {
	cases := []struct {
		Raw  string
		Type TokenType
	}{
		{
			Raw:  "fn",
			Type: KeyFn,
		},
		{
			Raw:  "print",
			Type: KeyPrint,
		},
		{
			Raw:  "str",
			Type: KeyStr,
		},
		{
			Raw:  "(",
			Type: OpLeftParen,
		},
		{
			Raw:  ")",
			Type: OpRightParen,
		},
		{
			Raw:  "{",
			Type: OpLeftBracket,
		},
		{
			Raw:  "}",
			Type: OpRightBracket,
		},
		{
			Raw:  "123",
			Type: TypeInt,
		},
		{
			Raw:  "\"hello\"",
			Type: TypeString,
		},
		{
			Raw:  "long_name_123",
			Type: TypeName,
		},
	}

	for _, c := range cases {
		t.Run(c.Raw, testToken(c.Raw, c.Type))
	}
}

func testToken(raw string, typ TokenType) func(*testing.T) {
	return func(t *testing.T) {
		token := NewToken(raw)

		if token.String() != raw {
			t.Errorf("expected %s, got %s", raw, token.String())
		}

		if token.Type() != typ {
			t.Errorf("expected %s, got %s", typ, token.Type())
		}
	}
}
