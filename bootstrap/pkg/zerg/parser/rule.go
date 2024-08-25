package parser

import (
	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/token"
)

// The LL(1) parser for the zerg language that generate the AST from the
// processed tokens.
type Rule interface {
	// Handle the processed token and build the AST by the self-hosed rules.
	//
	// It may return the unused token for the next rule to handle or just return nil if all
	// the tokens are consumed. It may consumed more than 2 tokens and acquire the next token
	// from the holder channel.
	Handle(ast *AST, prev, token *token.Token, holder <-chan *token.Token) (*token.Token, error)
}
