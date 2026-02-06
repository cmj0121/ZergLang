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
	// mut x := expr
	if p.curToken.Type == lexer.MUT {
		return p.parseMutableDeclaration()
	}

	// x := expr
	if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.DECLARE {
		return p.parseDeclarationStatement(false)
	}

	// x = expr or x, y = a, b
	if p.curToken.Type == lexer.IDENT && (p.peekToken.Type == lexer.ASSIGN || p.peekToken.Type == lexer.COMMA) {
		return p.parseAssignmentStatement()
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseMutableDeclaration() *DeclarationStatement {
	p.nextToken() // skip 'mut'

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("expected identifier after 'mut', got %s", p.curToken.Type))
		return nil
	}

	return p.parseDeclarationStatement(true)
}

func (p *Parser) parseDeclarationStatement(mutable bool) *DeclarationStatement {
	stmt := &DeclarationStatement{Mutable: mutable}

	name := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt.Name = name

	p.nextToken() // move to :=
	stmt.Token = p.curToken

	p.nextToken() // move to expression
	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseAssignmentStatement() *AssignmentStatement {
	stmt := &AssignmentStatement{}

	// Collect left-hand side identifiers
	stmt.Names = append(stmt.Names, &Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // move to comma
		p.nextToken() // move to next identifier

		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("expected identifier, got %s", p.curToken.Type))
			return nil
		}
		stmt.Names = append(stmt.Names, &Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if p.peekToken.Type != lexer.ASSIGN {
		p.errors = append(p.errors, fmt.Sprintf("expected '=', got %s", p.peekToken.Type))
		return nil
	}

	p.nextToken() // move to =
	stmt.Token = p.curToken

	// Collect right-hand side expressions
	p.nextToken() // move to first expression
	expr := p.parseExpression()
	if expr != nil {
		stmt.Values = append(stmt.Values, expr)
	}

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // move to comma
		p.nextToken() // move to next expression

		expr := p.parseExpression()
		if expr != nil {
			stmt.Values = append(stmt.Values, expr)
		}
	}

	if len(stmt.Names) != len(stmt.Values) {
		p.errors = append(p.errors, fmt.Sprintf("assignment count mismatch: %d names, %d values",
			len(stmt.Names), len(stmt.Values)))
		return nil
	}

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
