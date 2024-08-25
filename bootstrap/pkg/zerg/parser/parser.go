package parser

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/lexer"
	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/token"
)

// The hand-made parser for the zerg language that generate the AST from the
// processed tokens.
type Parser struct {
	*lexer.Lexer

	ast   *AST
	rules map[token.TokenType]Rule
}

// Create a new instance of the parser that holds the lexer and the parser rules.
func New(r io.Reader) *Parser {
	return &Parser{
		Lexer: lexer.New(r),
	}
}

// Parse from the tokenlized source code to the parsed AST.
func (p *Parser) Parse(ctx context.Context) error {
	p.prologue()
	defer p.epilogue()

	return p.parse(ctx)
}

// setup everything before the parsing, like register the parser rules.
func (p *Parser) prologue() {
	log.Debug().Msg("starting the parsing ...")
}

// clean up and release the resources after the parsing.
func (p *Parser) epilogue() {
	log.Debug().Msg("finished the parsing ...")
}

// start parsing the source code to the AST.
func (p *Parser) parse(ctx context.Context) error {
	holder := p.Lexer.Iterate(ctx)

	prev := <-holder

	// iterate the tokens until the EOF
	for prev != nil && prev.Type() != token.EndOfFile {
		var err error

		select {
		case <-ctx.Done():
			log.Warn().Msg("context is canceled")
			return p.Lexer.Err()
		case token := <-holder:
			// may acquired the nil token from the holder and only process the prev token
			switch rule, ok := p.rules[prev.Type()]; ok {
			case false:
				log.Warn().Str("token", prev.String()).Msg("no rule to handle the token")
				return fmt.Errorf("no rule to handle the token: %v", prev)
			case true:
				if prev, err = rule.Handle(p.ast, prev, token, holder); err != nil {
					log.Warn().Err(err).Str("token", token.String()).Msg("failed to handle the token")
					return err
				}

				if prev == nil {
					log.Debug().Msg("all tokens are consumed, acquring the next token")
					prev = <-holder
				}
			}
		}
	}

	return p.Lexer.Err()
}
