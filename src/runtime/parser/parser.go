package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xrspace/zerglang/runtime/lexer"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	OR_PREC      // or
	AND_PREC     // and
	EQUALS       // == !=
	LESSGREATER  // < > <= >=
	SUM          // + -
	PRODUCT      // * / %
	POWER_PREC   // **
	PREFIX       // -x, not x
	CALL         // fn()
)

var precedences = map[lexer.TokenType]int{
	lexer.OR:       OR_PREC,
	lexer.AND:      AND_PREC,
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.LT_EQ:    LESSGREATER,
	lexer.GT_EQ:    LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTERISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
	lexer.PERCENT:  PRODUCT,
	lexer.POWER:    POWER_PREC,
	lexer.LPAREN:   CALL,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// Parser produces an AST from a sequence of tokens.
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

// New creates a new Parser for the given Lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBoolean)
	p.registerPrefix(lexer.FALSE, p.parseBoolean)
	p.registerPrefix(lexer.NIL, p.parseNil)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.FN, p.parseFunctionLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.PERCENT, p.parseInfixExpression)
	p.registerInfix(lexer.POWER, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.GT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) parseStatement() Statement {
	// return expr
	if p.curToken.Type == lexer.RETURN {
		return p.parseReturnStatement()
	}

	// fn name(...) { ... } - named function declaration
	if p.curToken.Type == lexer.FN && p.peekToken.Type == lexer.IDENT {
		return p.parseFunctionDeclaration()
	}

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
	stmt.Value = p.parseExpression(LOWEST)

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
	expr := p.parseExpression(LOWEST)
	if expr != nil {
		stmt.Values = append(stmt.Values, expr)
	}

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // move to comma
		p.nextToken() // move to next expression

		expr := p.parseExpression(LOWEST)
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
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
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

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == lexer.TRUE}
}

func (p *Parser) parseNil() Expression {
	return &NilLiteral{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()

	// Right-associative for power operator
	if p.curToken.Type == lexer.POWER {
		precedence--
	}

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, "expected )")
		return nil
	}
	p.nextToken()

	return exp
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	if p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	return stmt
}

func (p *Parser) parseFunctionDeclaration() *DeclarationStatement {
	fnToken := p.curToken
	p.nextToken() // move to function name

	name := &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	fn := p.parseFunctionLiteralWithName(fnToken, name)

	return &DeclarationStatement{
		Token:   fnToken,
		Name:    name,
		Value:   fn,
		Mutable: false,
	}
}

func (p *Parser) parseFunctionLiteral() Expression {
	return p.parseFunctionLiteralWithName(p.curToken, nil)
}

func (p *Parser) parseFunctionLiteralWithName(fnToken lexer.Token, name *Identifier) *FunctionLiteral {
	lit := &FunctionLiteral{Token: fnToken, Name: name}

	if p.peekToken.Type != lexer.LPAREN {
		p.errors = append(p.errors, "expected ( after fn")
		return nil
	}
	p.nextToken()

	lit.Parameters = p.parseFunctionParameters()

	// Skip optional return type annotation: -> type
	if p.peekToken.Type == lexer.ARROW {
		p.nextToken() // move to ->
		p.nextToken() // move to type, skip it
	}

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, "expected { for function body")
		return nil
	}
	p.nextToken()

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*Parameter {
	params := []*Parameter{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := p.parseParameter()
	params = append(params, param)

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // move to comma
		p.nextToken() // move to next param
		param := p.parseParameter()
		params = append(params, param)
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, "expected ) after parameters")
		return nil
	}
	p.nextToken()

	return params
}

func (p *Parser) parseParameter() *Parameter {
	param := &Parameter{
		Name: &Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	// Skip type annotation: name: type
	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // move to :
		p.nextToken() // move to type, skip it
	}

	// Parse default value: name = expr
	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken() // move to =
		p.nextToken() // move to default value
		param.Default = p.parseExpression(LOWEST)
	}

	return param
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, "expected ) after arguments")
		return nil
	}
	p.nextToken()

	return args
}
