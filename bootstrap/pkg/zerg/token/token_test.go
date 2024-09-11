package token

import (
	"testing"
)

func TestToken(t *testing.T) {
	cases := []struct {
		Raw  string
		Type Type
	}{
		{
			Raw:  "fn",
			Type: Fn,
		},
		{
			Raw:  "print",
			Type: Print,
		},
		{
			Raw:  "str",
			Type: Str,
		},
		{
			Raw:  "(",
			Type: LParen,
		},
		{
			Raw:  ")",
			Type: RParen,
		},
		{
			Raw:  "{",
			Type: LBrace,
		},
		{
			Raw:  "}",
			Type: RBrace,
		},
		{
			Raw:  "123",
			Type: Int,
		},
		{
			Raw:  "\"hello\"",
			Type: String,
		},
		{
			Raw:  "long_name_123",
			Type: Name,
		},
	}

	for _, c := range cases {
		t.Run(c.Raw, testToken(c.Raw, c.Type))
	}
}

func testToken(raw string, typ Type) func(*testing.T) {
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
