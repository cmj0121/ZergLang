package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xrspace/zerglang/runtime/lexer"
)

// Parser produces an AST from a sequence of tokens.
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

// New creates a new Parser for the given Lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

// ParseProgram parses the entire program and returns the AST.
func (p *Parser) ParseProgram() *Program {
	program := &Program{Statements: []Statement{}}

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// Errors returns any parsing errors encountered.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() Statement {
	if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.DECLARE {
		return p.parseDeclarationStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseDeclarationStatement() *DeclarationStatement {
	stmt := &DeclarationStatement{}

	name := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt.Name = name

	p.nextToken() // move to :=
	stmt.Token = p.curToken

	p.nextToken() // move to expression
	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression()
	return stmt
}

func (p *Parser) parseExpression() Expression {
	switch p.curToken.Type {
	case lexer.IDENT:
		return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case lexer.INT:
		return p.parseIntegerLiteral()
	case lexer.STRING:
		return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case lexer.TRUE:
		return &BooleanLiteral{Token: p.curToken, Value: true}
	case lexer.FALSE:
		return &BooleanLiteral{Token: p.curToken, Value: false}
	case lexer.NIL:
		return &NilLiteral{Token: p.curToken}
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s", p.curToken.Type))
		return nil
	}
}

func (p *Parser) parseIntegerLiteral() *IntegerLiteral {
	lit := &IntegerLiteral{Token: p.curToken}

	// Remove underscores from the literal
	cleanLiteral := strings.ReplaceAll(p.curToken.Literal, "_", "")

	value, err := strconv.ParseInt(cleanLiteral, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}
