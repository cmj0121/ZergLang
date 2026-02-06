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
	INDEX        // arr[0], obj.field
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
	lexer.LBRACKET: INDEX,
	lexer.DOT:      INDEX,
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
	p.registerPrefix(lexer.LBRACKET, p.parseListLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseMapLiteral)
	p.registerPrefix(lexer.THIS, p.parseThis)
	p.registerPrefix(lexer.SELF, p.parseSelf)

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
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.DOT, p.parseMemberExpression)

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
	switch p.curToken.Type {
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.BREAK:
		return &BreakStatement{Token: p.curToken}
	case lexer.CONTINUE:
		return &ContinueStatement{Token: p.curToken}
	case lexer.NOP:
		return &NopStatement{Token: p.curToken}
	case lexer.CLASS:
		return p.parseClassDeclaration()
	case lexer.SPEC:
		return p.parseSpecDeclaration()
	case lexer.IMPL:
		return p.parseImplDeclaration()
	case lexer.FN:
		if p.peekToken.Type == lexer.IDENT {
			return p.parseFunctionDeclaration()
		}
	case lexer.MUT:
		return p.parseMutableDeclaration()
	case lexer.IDENT:
		if p.peekToken.Type == lexer.DECLARE {
			return p.parseDeclarationStatement(false)
		}
		if p.peekToken.Type == lexer.ASSIGN || p.peekToken.Type == lexer.COMMA {
			return p.parseAssignmentStatement()
		}
		// Check for member/index assignment: ident.field = value or ident[idx] = value
		if p.peekToken.Type == lexer.DOT || p.peekToken.Type == lexer.LBRACKET {
			return p.tryParseMemberAssignment()
		}
	case lexer.THIS:
		// Check for member assignment: this.field = value
		if p.peekToken.Type == lexer.DOT {
			return p.tryParseMemberAssignment()
		}
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

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, "expected { after if condition")
		return nil
	}
	p.nextToken()
	stmt.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == lexer.ELSE {
		p.nextToken()

		if p.peekToken.Type != lexer.LBRACE {
			p.errors = append(p.errors, "expected { after else")
			return nil
		}
		p.nextToken()
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseForStatement() Statement {
	token := p.curToken

	// for { } - infinite loop
	if p.peekToken.Type == lexer.LBRACE {
		p.nextToken()
		body := p.parseBlockStatement()
		return &ForConditionStatement{Token: token, Condition: nil, Body: body}
	}

	p.nextToken()

	// Check if it's a for-in loop: for item in collection { }
	if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.IN {
		variable := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken() // move to 'in'
		p.nextToken() // move to iterable

		iterable := p.parseExpression(LOWEST)

		if p.peekToken.Type != lexer.LBRACE {
			p.errors = append(p.errors, "expected { after for-in")
			return nil
		}
		p.nextToken()
		body := p.parseBlockStatement()

		return &ForInStatement{Token: token, Variable: variable, Iterable: iterable, Body: body}
	}

	// for condition { } - conditional loop
	condition := p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, "expected { after for condition")
		return nil
	}
	p.nextToken()
	body := p.parseBlockStatement()

	return &ForConditionStatement{Token: token, Condition: condition, Body: body}
}

func (p *Parser) parseListLiteral() Expression {
	list := &ListLiteral{Token: p.curToken}
	list.Elements = p.parseExpressionList(lexer.RBRACKET)
	return list
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []Expression {
	list := []Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if p.peekToken.Type != end {
		p.errors = append(p.errors, fmt.Sprintf("expected %s", end))
		return nil
	}
	p.nextToken()

	return list
}

func (p *Parser) parseMapLiteral() Expression {
	mapLit := &MapLiteral{Token: p.curToken}
	mapLit.Pairs = make(map[Expression]Expression)

	for p.peekToken.Type != lexer.RBRACE {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if p.peekToken.Type != lexer.COLON {
			p.errors = append(p.errors, "expected : in map literal")
			return nil
		}
		p.nextToken()

		p.nextToken()
		value := p.parseExpression(LOWEST)

		mapLit.Pairs[key] = value

		if p.peekToken.Type != lexer.RBRACE && p.peekToken.Type != lexer.COMMA {
			p.errors = append(p.errors, "expected } or , in map literal")
			return nil
		}

		if p.peekToken.Type == lexer.COMMA {
			p.nextToken()
		}
	}

	if p.peekToken.Type != lexer.RBRACE {
		p.errors = append(p.errors, "expected } at end of map literal")
		return nil
	}
	p.nextToken()

	return mapLit
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.RBRACKET {
		p.errors = append(p.errors, "expected ] after index")
		return nil
	}
	p.nextToken()

	return exp
}

func (p *Parser) parseMemberExpression(left Expression) Expression {
	exp := &MemberExpression{Token: p.curToken, Object: left}

	p.nextToken()
	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "expected identifier after .")
		return nil
	}

	exp.Member = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return exp
}

func (p *Parser) parseThis() Expression {
	return &ThisExpression{Token: p.curToken}
}

func (p *Parser) tryParseMemberAssignment() Statement {
	// Parse the left-hand side expression first
	expr := p.parseExpression(LOWEST)

	// Check if followed by = for assignment
	if p.peekToken.Type != lexer.ASSIGN {
		// Not an assignment, return as expression statement
		return &ExpressionStatement{Token: p.curToken, Expression: expr}
	}

	// It's a member/index assignment
	switch left := expr.(type) {
	case *MemberExpression:
		p.nextToken() // move to =
		assignToken := p.curToken
		p.nextToken() // move to value
		value := p.parseExpression(LOWEST)
		return &MemberAssignmentStatement{
			Token:  assignToken,
			Object: left.Object,
			Member: left.Member,
			Value:  value,
		}
	case *IndexExpression:
		p.nextToken() // move to =
		assignToken := p.curToken
		p.nextToken() // move to value
		value := p.parseExpression(LOWEST)
		return &IndexAssignmentStatement{
			Token: assignToken,
			Left:  left.Left,
			Index: left.Index,
			Value: value,
		}
	default:
		// Return as expression statement
		return &ExpressionStatement{Token: p.curToken, Expression: expr}
	}
}

func (p *Parser) parseClassDeclaration() *ClassDeclaration {
	stmt := &ClassDeclaration{Token: p.curToken}

	if p.peekToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "expected class name")
		return nil
	}
	p.nextToken()
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, "expected { after class name")
		return nil
	}
	p.nextToken()

	stmt.Fields = p.parseClassFields()

	return stmt
}

func (p *Parser) parseClassFields() []*FieldDeclaration {
	fields := []*FieldDeclaration{}

	p.nextToken()

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		field := p.parseFieldDeclaration()
		if field != nil {
			fields = append(fields, field)
		}
		p.nextToken()
	}

	return fields
}

func (p *Parser) parseFieldDeclaration() *FieldDeclaration {
	field := &FieldDeclaration{}

	// Check for pub modifier
	if p.curToken.Type == lexer.PUB {
		field.Public = true
		p.nextToken()
	}

	// Check for mut modifier
	if p.curToken.Type == lexer.MUT {
		field.Mutable = true
		p.nextToken()
	}

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "expected field name")
		return nil
	}

	field.Token = p.curToken
	field.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Skip type annotation if present: name: type
	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // move to :
		p.nextToken() // move to type, skip it
	}

	// Parse default value if present: name = expr
	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken() // move to =
		p.nextToken() // move to default value
		field.Default = p.parseExpression(LOWEST)
	}

	return field
}

func (p *Parser) parseImplDeclaration() Statement {
	implToken := p.curToken

	if p.peekToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "expected class name after impl")
		return nil
	}
	p.nextToken()
	className := &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for "impl Class for Spec" syntax
	if p.peekToken.Type == lexer.FOR {
		p.nextToken() // move to 'for'

		if p.peekToken.Type != lexer.IDENT {
			p.errors = append(p.errors, "expected spec name after 'for'")
			return nil
		}
		p.nextToken()
		specName := &Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.peekToken.Type != lexer.LBRACE {
			p.errors = append(p.errors, "expected { after spec name")
			return nil
		}
		p.nextToken()

		methods := p.parseMethodDeclarations()

		return &ImplForDeclaration{
			Token:   implToken,
			Class:   className,
			Spec:    specName,
			Methods: methods,
		}
	}

	// Regular "impl Class" syntax
	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, "expected { after class name")
		return nil
	}
	p.nextToken()

	methods := p.parseMethodDeclarations()

	return &ImplDeclaration{
		Token:   implToken,
		Class:   className,
		Methods: methods,
	}
}

func (p *Parser) parseMethodDeclarations() []*MethodDeclaration {
	methods := []*MethodDeclaration{}

	p.nextToken()

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		method := p.parseMethodDeclaration()
		if method != nil {
			methods = append(methods, method)
		}
		p.nextToken()
	}

	return methods
}

func (p *Parser) parseMethodDeclaration() *MethodDeclaration {
	method := &MethodDeclaration{}

	// Check for pub modifier
	if p.curToken.Type == lexer.PUB {
		method.Public = true
		p.nextToken()
	}

	// Check for static modifier
	if p.curToken.Type == lexer.STATIC {
		method.Static = true
		p.nextToken()
	}

	// Check for mut modifier (mutable receiver)
	if p.curToken.Type == lexer.MUT {
		method.Mutable = true
		p.nextToken()
	}

	if p.curToken.Type != lexer.FN {
		p.errors = append(p.errors, "expected fn keyword in method")
		return nil
	}
	method.Token = p.curToken

	if p.peekToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "expected method name")
		return nil
	}
	p.nextToken()
	method.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != lexer.LPAREN {
		p.errors = append(p.errors, "expected ( after method name")
		return nil
	}
	p.nextToken()

	method.Parameters = p.parseFunctionParameters()

	// Skip optional return type annotation: -> type
	if p.peekToken.Type == lexer.ARROW {
		p.nextToken() // move to ->
		p.nextToken() // move to type, skip it
	}

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, "expected { for method body")
		return nil
	}
	p.nextToken()

	method.Body = p.parseBlockStatement()

	return method
}

func (p *Parser) parseSpecDeclaration() *SpecDeclaration {
	stmt := &SpecDeclaration{Token: p.curToken}

	if p.peekToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "expected spec name")
		return nil
	}
	p.nextToken()
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, "expected { after spec name")
		return nil
	}
	p.nextToken()

	stmt.Methods = p.parseMethodSignatures()

	return stmt
}

func (p *Parser) parseMethodSignatures() []*MethodSignature {
	methods := []*MethodSignature{}

	p.nextToken()

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		method := p.parseMethodSignature()
		if method != nil {
			methods = append(methods, method)
		}
		p.nextToken()
	}

	return methods
}

func (p *Parser) parseMethodSignature() *MethodSignature {
	method := &MethodSignature{}

	// Check for pub modifier
	if p.curToken.Type == lexer.PUB {
		method.Public = true
		p.nextToken()
	}

	// Check for mut modifier (mutable receiver)
	if p.curToken.Type == lexer.MUT {
		method.Mutable = true
		p.nextToken()
	}

	if p.curToken.Type != lexer.FN {
		p.errors = append(p.errors, "expected fn keyword in method signature")
		return nil
	}
	method.Token = p.curToken

	if p.peekToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "expected method name")
		return nil
	}
	p.nextToken()
	method.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != lexer.LPAREN {
		p.errors = append(p.errors, "expected ( after method name")
		return nil
	}
	p.nextToken()

	// Parse parameter names only (no types in bootstrap)
	method.Parameters = p.parseSignatureParameters()

	// Skip optional return type annotation: -> type
	if p.peekToken.Type == lexer.ARROW {
		p.nextToken() // move to ->
		p.nextToken() // move to type, skip it
	}

	return method
}

func (p *Parser) parseSignatureParameters() []*Identifier {
	params := []*Identifier{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
		return params
	}

	p.nextToken()
	params = append(params, &Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // move to comma
		p.nextToken() // move to next param
		params = append(params, &Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, "expected ) after parameters")
		return nil
	}
	p.nextToken()

	return params
}

func (p *Parser) parseSelf() Expression {
	return &SelfExpression{Token: p.curToken}
}
