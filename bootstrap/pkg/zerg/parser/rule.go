package parser

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/token"
)

// Handle the processed token and build the AST by the self-hosed rules.
//
// It may return the unused token for the next rule to handle or just return nil if all
// the tokens are consumed. It may consumed more than 2 tokens and acquire the next token
// from the holder channel.
type Rule func(root *Node, prev, token *token.Token, holder <-chan *token.Token) (*token.Token, error)

// Parse the function statement.
//
// func ::= "fn" NAME func_args scope
func RuleFunc(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Any("prev", prev.Type()).Any("curr", curr.Type()).Msg("parse the function")

	if prev.Type() != token.Fn {
		return nil, fmt.Errorf("expect the function declaration but got %s", prev)
	}

	if curr.Type() != token.Name {
		return nil, fmt.Errorf("expect the function name but got %s", curr)
	}

	var node = &Node{typ: FN, token: curr}
	var err error

	root.Append(node)

	prev = <-holder
	curr = <-holder

	switch prev, err = RuleFuncArgs(node, prev, curr, holder); err {
	case nil:
		switch prev {
		case nil:
			prev, curr = <-holder, <-holder
		default:
			curr = <-holder
		}
	default:
		log.Debug().Err(err).Msg("failed to parse the function arguments")
		return nil, err
	}

	return RuleScope(node, prev, curr, holder)
}

// Parse the function arguments.
//
// func_args ::= "(" ")"
func RuleFuncArgs(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Any("prev", prev.Type()).Any("curr", curr.Type()).Msg("parse the function arguments")

	if prev.Type() != token.LParen {
		return nil, fmt.Errorf("expect the left parenthesis but got %s", prev)
	}

	if curr.Type() != token.RParen {
		return nil, fmt.Errorf("expect the right parenthesis but got %s", curr)
	}

	args := &Node{typ: ARGS}
	root.Append(args)

	return nil, nil
}

// Parse the scope that contains the statements.
//
// scope ::= "{" stmt+ "}"
func RuleScope(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Any("prev", prev.Type()).Any("curr", curr.Type()).Msg("parse the scope")

	if prev.Type() != token.LBrace {
		return nil, fmt.Errorf("expect the left bracket but got %s", prev)
	}

	prev, curr = curr, <-holder
	if prev.Type() == token.RBrace {
		// the final right Brace
		return curr, nil
	}

	var scope = &Node{typ: SCOPE}
	root.Append(scope)

	err := fmt.Errorf("not implemented: %v", prev)
	return nil, err
}
