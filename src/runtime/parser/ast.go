package parser

import "github.com/xrspace/zerglang/runtime/lexer"

// Node is the base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
}

// Statement represents a statement node.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// DeclarationStatement represents a variable declaration: x := expr or mut x := expr
type DeclarationStatement struct {
	Token   lexer.Token // the := token
	Name    *Identifier
	Value   Expression
	Mutable bool
}

func (ds *DeclarationStatement) statementNode()       {}
func (ds *DeclarationStatement) TokenLiteral() string { return ds.Token.Literal }

// AssignmentStatement represents variable reassignment: x = expr or x, y = a, b
type AssignmentStatement struct {
	Token  lexer.Token  // the = token
	Names  []*Identifier
	Values []Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }

// ExpressionStatement represents a bare expression used as a statement.
type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// Identifier represents a variable name.
type Identifier struct {
	Token lexer.Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// IntegerLiteral represents an integer value.
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// StringLiteral represents a string value.
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// BooleanLiteral represents a boolean value.
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }

// NilLiteral represents the nil value.
type NilLiteral struct {
	Token lexer.Token
}

func (nl *NilLiteral) expressionNode()      {}
func (nl *NilLiteral) TokenLiteral() string { return nl.Token.Literal }

// PrefixExpression represents a prefix operator: -x, not x
type PrefixExpression struct {
	Token    lexer.Token // the prefix token, e.g. - or not
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// InfixExpression represents an infix operator: x + y, x and y
type InfixExpression struct {
	Token    lexer.Token // the operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// BlockStatement represents a block of statements: { ... }
type BlockStatement struct {
	Token      lexer.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// Parameter represents a function parameter.
type Parameter struct {
	Name    *Identifier
	Default Expression // optional default value
}

// FunctionLiteral represents a function: fn name(params) -> type { body }
type FunctionLiteral struct {
	Token      lexer.Token // the 'fn' token
	Name       *Identifier // nil for anonymous functions
	Parameters []*Parameter
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// CallExpression represents a function call: fn(args)
type CallExpression struct {
	Token     lexer.Token // the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// ReturnStatement represents a return statement.
type ReturnStatement struct {
	Token       lexer.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
