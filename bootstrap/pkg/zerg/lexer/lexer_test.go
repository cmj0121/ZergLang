package lexer

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func Example() {
	source := `
		fn main() {
			print "Hello, World!"
		}
	`

	reader := strings.NewReader(source)
	lexer := New(reader)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for token := range lexer.Iterate(ctx) {
		fmt.Println(token.String())
	}

	// Output:
	// fn
	// main
	// (
	// )
	// {
	// print
	// "Hello, World!"
	// }
}

func TestLexer(t *testing.T) {
	cases := []struct {
		Source string
		Tokens []string
	}{
		{
			Source: "fn main(){}",
			Tokens: []string{"fn", "main", "(", ")", "{", "}"},
		},
		{
			Source: "long_name_123",
			Tokens: []string{"long_name_123"},
		},
		{
			Source: `"The text with emoji 🤣"`,
			Tokens: []string{`"The text with emoji 🤣"`},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case-%d", i), testLexer(c.Source, c.Tokens))
	}
}

func testLexer(source string, tokens []string) func(*testing.T) {
	return func(t *testing.T) {
		reader := strings.NewReader(source)
		lexer := New(reader)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		index := 0
		for token := range lexer.Iterate(ctx) {
			if len(tokens) == 0 {
				t.Fatalf("unexpected token %q", token.String())
			}

			expected_token := tokens[index]
			tokens = tokens[1:]

			if token.String() != expected_token {
				t.Errorf("expect %q, got %q", expected_token, token.String())
			}
		}

		if len(tokens) > 0 {
			t.Errorf("expect %d more tokens", len(tokens))
		}
	}
}
