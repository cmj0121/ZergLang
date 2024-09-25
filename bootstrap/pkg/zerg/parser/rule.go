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
	log.Debug().Str("prev", prev.String()).Msg("parse the function")

	if prev.Type() != token.Fn {
		return nil, fmt.Errorf("expect the function declaration but got %s", prev)
	}

	if curr.Type() != token.Name {
		return nil, fmt.Errorf("expect the function name but got %s", curr)
	}

	var node = &Node{typ: Fn, token: curr}
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
// func_args ::= "(" ")" type_hint?
func RuleFuncArgs(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Str("prev", prev.String()).Msg("parse the function arguments")

	if prev.Type() != token.LParen {
		return nil, fmt.Errorf("expect the left parenthesis but got %s", prev)
	}

	if curr.Type() != token.RParen {
		return nil, fmt.Errorf("expect the right parenthesis but got %s", curr)
	}

	args := &Node{typ: Args}
	root.Append(args)

	switch prev = <-holder; prev.Type() {
	case token.Arrow:
		curr = <-holder
		return RuleTypeHint(root, prev, curr, holder)
	default:
		type_hint := &Node{typ: Type}
		root.Append(type_hint)
		return prev, nil
	}
}

// Parse the scope that contains the statements.
//
// scope ::= "{" stmt+ "}"
func RuleScope(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Str("prev", prev.String()).Msg("parse the scope")

	if prev.Type() != token.LBrace {
		return nil, fmt.Errorf("expect the left bracket but got %s", prev)
	}

	var err error
	var scope = &Node{typ: Scope}
	root.Append(scope)

	for {
		prev, curr = curr, <-holder

		if prev.Type() == token.RBrace {
			// the final right Brace
			return curr, nil
		}

		curr, err = RuleStmt(scope, prev, curr, holder)
		if err != nil {
			log.Info().Err(err).Msg("failed to parse the statement")
			return nil, err
		}
	}
}

// Parse the type hint.
//
// type_hint ::= "->" Type
func RuleTypeHint(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Str("prev", prev.String()).Msg("parse the type hint")

	if prev.Type() != token.Arrow {
		return nil, fmt.Errorf("expect the arrow but got %s", prev)
	}

	switch curr.Type() {
	case token.Name:
		var node = &Node{typ: Type, token: curr}
		root.Append(node)
	default:
		return nil, fmt.Errorf("expect the type name but got %s", curr)
	}

	return nil, nil
}

// Parse the statement.
//
// stmt ::= return_stmt
func RuleStmt(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Str("prev", prev.String()).Msg("parse the statement")

	switch prev.Type() {
	case token.Return:
		return RuleReturn(root, prev, curr, holder)
	case token.Print:
		return RulePrint(root, prev, curr, holder)
	default:
		return nil, fmt.Errorf("unknown statement: %s", prev)
	}
}

// Parse the return statement.
//
// return_stmt ::= "return" expr
func RuleReturn(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Str("prev", prev.String()).Msg("parse the return statement")

	if prev.Type() != token.Return {
		return nil, fmt.Errorf("expect the return statement but got %s", prev)
	}

	var node = &Node{typ: ReturnStmt, token: prev}
	root.Append(node)

	prev, curr = curr, <-holder
	return RuleExpr(node, prev, curr, holder)
}

// Parse the print statement.
//
// print_stmt ::= "print" expr
func RulePrint(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Str("prev", prev.String()).Msg("parse the print statement")

	if prev.Type() != token.Print {
		return nil, fmt.Errorf("expect the print statement but got %s", prev)
	}

	var node = &Node{typ: PrintStmt, token: prev}
	root.Append(node)

	prev, curr = curr, <-holder
	return RuleExpr(node, prev, curr, holder)
}

// Parse the expression.
//
// expr ::= NAME
func RuleExpr(root *Node, prev, curr *token.Token, holder <-chan *token.Token) (*token.Token, error) {
	log.Debug().Str("prev", prev.String()).Msg("parse the expression")

	switch prev.Type() {
	case token.Name, token.String, token.Int:
		var node = &Node{typ: Expression, token: prev}
		root.Append(node)
	default:
		return nil, fmt.Errorf("expect the expression but got %s", prev)
	}

	return curr, nil
}
