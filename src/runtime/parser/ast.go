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

// IfStatement represents an if statement.
type IfStatement struct {
	Token       lexer.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // optional
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }

// ForInStatement represents a for-in loop: for item in collection { }
type ForInStatement struct {
	Token      lexer.Token // the 'for' token
	Variable   *Identifier
	Iterable   Expression
	Body       *BlockStatement
}

func (fis *ForInStatement) statementNode()       {}
func (fis *ForInStatement) TokenLiteral() string { return fis.Token.Literal }

// ForConditionStatement represents a for loop with condition: for condition { }
type ForConditionStatement struct {
	Token     lexer.Token // the 'for' token
	Condition Expression  // nil for infinite loop
	Body      *BlockStatement
}

func (fcs *ForConditionStatement) statementNode()       {}
func (fcs *ForConditionStatement) TokenLiteral() string { return fcs.Token.Literal }

// BreakStatement represents a break statement.
type BreakStatement struct {
	Token lexer.Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }

// ContinueStatement represents a continue statement.
type ContinueStatement struct {
	Token lexer.Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }

// NopStatement represents a no-operation statement.
type NopStatement struct {
	Token lexer.Token
}

func (ns *NopStatement) statementNode()       {}
func (ns *NopStatement) TokenLiteral() string { return ns.Token.Literal }

// ListLiteral represents a list: [1, 2, 3]
type ListLiteral struct {
	Token    lexer.Token // the '[' token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }

// MapLiteral represents a map: {key: value, ...}
type MapLiteral struct {
	Token lexer.Token // the '{' token
	Pairs map[Expression]Expression
}

func (ml *MapLiteral) expressionNode()      {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }

// IndexExpression represents index access: arr[0], map["key"]
type IndexExpression struct {
	Token lexer.Token // the '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// MemberExpression represents member access: obj.field
type MemberExpression struct {
	Token  lexer.Token // the '.' token
	Object Expression
	Member *Identifier
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }

// ClassDeclaration represents a class definition: class Name { fields }
type ClassDeclaration struct {
	Token  lexer.Token // the 'class' token
	Name   *Identifier
	Fields []*FieldDeclaration
}

func (cd *ClassDeclaration) statementNode()       {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }

// FieldDeclaration represents a class field: pub mut name: type = default
type FieldDeclaration struct {
	Token   lexer.Token // the field name token
	Name    *Identifier
	Public  bool
	Mutable bool
	Default Expression // optional default value
}

func (fd *FieldDeclaration) statementNode()       {}
func (fd *FieldDeclaration) TokenLiteral() string { return fd.Token.Literal }

// ImplDeclaration represents method implementations: impl ClassName { methods }
type ImplDeclaration struct {
	Token   lexer.Token // the 'impl' token
	Class   *Identifier
	Methods []*MethodDeclaration
}

func (id *ImplDeclaration) statementNode()       {}
func (id *ImplDeclaration) TokenLiteral() string { return id.Token.Literal }

// MethodDeclaration represents a method: pub static mut fn name(params) { body }
type MethodDeclaration struct {
	Token      lexer.Token // the 'fn' token
	Name       *Identifier
	Parameters []*Parameter
	Body       *BlockStatement
	Public     bool
	Static     bool
	Mutable    bool // mut receiver (self)
}

func (md *MethodDeclaration) statementNode()       {}
func (md *MethodDeclaration) TokenLiteral() string { return md.Token.Literal }

// ThisExpression represents 'this' keyword.
type ThisExpression struct {
	Token lexer.Token
}

func (te *ThisExpression) expressionNode()      {}
func (te *ThisExpression) TokenLiteral() string { return te.Token.Literal }

// MemberAssignmentStatement represents member assignment: obj.field = value
type MemberAssignmentStatement struct {
	Token  lexer.Token // the '=' token
	Object Expression
	Member *Identifier
	Value  Expression
}

func (mas *MemberAssignmentStatement) statementNode()       {}
func (mas *MemberAssignmentStatement) TokenLiteral() string { return mas.Token.Literal }

// IndexAssignmentStatement represents index assignment: arr[idx] = value
type IndexAssignmentStatement struct {
	Token lexer.Token // the '=' token
	Left  Expression
	Index Expression
	Value Expression
}

func (ias *IndexAssignmentStatement) statementNode()       {}
func (ias *IndexAssignmentStatement) TokenLiteral() string { return ias.Token.Literal }

// SpecDeclaration represents a spec (interface) definition: spec Name { methods }
type SpecDeclaration struct {
	Token   lexer.Token // the 'spec' token
	Name    *Identifier
	Methods []*MethodSignature
}

func (sd *SpecDeclaration) statementNode()       {}
func (sd *SpecDeclaration) TokenLiteral() string { return sd.Token.Literal }

// MethodSignature represents a method signature in a spec: fn name(params) -> type
type MethodSignature struct {
	Token      lexer.Token // the 'fn' token
	Name       *Identifier
	Parameters []*Identifier // just names for signatures
	Public     bool
	Mutable    bool // mut receiver
}

func (ms *MethodSignature) statementNode()       {}
func (ms *MethodSignature) TokenLiteral() string { return ms.Token.Literal }

// ImplForDeclaration represents impl Class for Spec { methods }
type ImplForDeclaration struct {
	Token   lexer.Token // the 'impl' token
	Class   *Identifier
	Spec    *Identifier
	Methods []*MethodDeclaration
}

func (ifd *ImplForDeclaration) statementNode()       {}
func (ifd *ImplForDeclaration) TokenLiteral() string { return ifd.Token.Literal }

// SelfExpression represents 'Self' type reference.
type SelfExpression struct {
	Token lexer.Token
}

func (se *SelfExpression) expressionNode()      {}
func (se *SelfExpression) TokenLiteral() string { return se.Token.Literal }

// ReferenceExpression represents a reference (&expr).
type ReferenceExpression struct {
	Token lexer.Token // The '&' token
	Value Expression
}

func (re *ReferenceExpression) expressionNode()      {}
func (re *ReferenceExpression) TokenLiteral() string { return re.Token.Literal }

// AssertStatement represents an assert statement.
type AssertStatement struct {
	Token     lexer.Token
	Condition Expression
	Message   Expression // optional message
}

func (as *AssertStatement) statementNode()       {}
func (as *AssertStatement) TokenLiteral() string { return as.Token.Literal }

// UnsafeBlock represents: unsafe { ... }
type UnsafeBlock struct {
	Token lexer.Token // the 'unsafe' token
	Body  *BlockStatement
}

func (ub *UnsafeBlock) statementNode()       {}
func (ub *UnsafeBlock) TokenLiteral() string { return ub.Token.Literal }

// AsmExpression represents: asm("go_function", args...)
type AsmExpression struct {
	Token    lexer.Token  // the 'asm' token
	Function string       // Go function name
	Args     []Expression // arguments
}

func (ae *AsmExpression) expressionNode()      {}
func (ae *AsmExpression) TokenLiteral() string { return ae.Token.Literal }
