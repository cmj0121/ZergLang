package lexer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/token"
)

// The custom error type for the tokenizer that show the line and position
// information that the error occurred.
type LexerError struct {
	err  error
	line int
	pos  int
}

func (e LexerError) Error() string {
	return fmt.Sprintf("line %d, pos %d: %s", e.line, e.pos, e.err)
}

// the source code tokenizer
type Lexer struct {
	r io.Reader

	line int
	buff string
	err  error
}

// create a new tokenizer
func New(r io.Reader) *Lexer {
	return &Lexer{
		r: r,
	}
}

// Get the last error that occurred during the tokenization
func (l *Lexer) Err() error {
	return l.err
}

// tokenize the source code that iterate the tokens and merge the
// multi-tokens into a single token if necessary.
func (l *Lexer) Iterate(ctx context.Context) <-chan *token.Token {
	ch := make(chan *token.Token)

	go func() {
		defer close(ch)

		prev := &token.EndOfLine
		for tt := range l.iterate(ctx) {
			switch prev.Type() {
			case token.EOL:
				prev = tt
				continue
			}

			select {
			case <-ctx.Done():
				return
			case ch <- prev:
				prev = tt
			}
		}

		// the last token should be EOF
		if prev.Type() != token.EOF {
			log.Warn().Str("token", prev.String()).Int("_line", l.line).Msg("the last token is not EOF")
		}
	}()

	return ch
}

// tokenize the source code and return the valid and processed tokens
func (l *Lexer) iterate(ctx context.Context) <-chan *token.Token {
	ch := make(chan *token.Token)

	go func() {
		defer close(ch)

		scanner := bufio.NewScanner(l.r)
		for scanner.Scan() {
			// read line into the instance buffer
			l.line++
			l.buff = strings.Trim(scanner.Text(), " \t")

			log.Trace().Str("text", l.buff).Int("line", l.line).Msg("read the line")
			// tokenize by line
			tokens, err := l.tokenize(l.buff)
			if err != nil {
				log.Warn().Err(err).Int("line", l.line).Msg("failed to tokenize the line")
				l.err = err
				return
			}

			// send the tokens to the channel
			for _, raw := range tokens {
				token := token.NewToken(raw)

				select {
				case <-ctx.Done():
					return
				case ch <- token:
					log.Debug().Str("token", token.String()).Interface("_typ", token.Type()).Int("_line", l.line).Msg("send the token")
				}
			}

			// return the EOL token
			ch <- &token.EndOfLine
		}

		// send the EOF token
		ch <- &token.EndOfFile
	}()

	return ch
}

// tokenlize by line with predefined rules
func (l *Lexer) tokenize(line string) ([]string, error) {
	tokens := []string{}

	index := 0
	for index < len(line) {
		switch ch := line[index]; ch {
		case ' ', '\t':
			// skip the white spaces
		case '/':
			if index+1 < len(line) && line[index+1] == '/' {
				// skip the comment
				return tokens, nil
			}

			tokens = append(tokens, line[index:index+1])
		case '"':
			found := false
			// find the string token
			for end := index + 1; end < len(line); end++ {
				if ch == line[end] {
					tokens = append(tokens, line[index:end+1])
					index = end
					found = true
					break
				}
			}

			if !found {
				err := LexerError{
					line: l.line,
					pos:  index,
					err:  fmt.Errorf("unterminated string"),
				}
				return nil, err
			}
		default:
			if !(unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_') {
				tokens = append(tokens, line[index:index+1])
				break
			}

			found := false
			// find the identifier token
			for end := index; end < len(line); end++ {
				if !(unicode.IsLetter(rune(line[end])) || unicode.IsDigit(rune(line[end])) || line[end] == '_') {
					tokens = append(tokens, line[index:end])
					index = end - 1
					found = true
					break
				}
			}

			if !found {
				tokens = append(tokens, line[index:])
				index = len(line)
			}
		}

		index++
	}

	return tokens, nil
}
