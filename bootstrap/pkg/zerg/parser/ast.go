package parser

import (
	"fmt"
	"strings"

	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/token"
)

// The type of the AST node.
//
//go:generate stringer -type=NodeType
type NodeType int

const (
	// The node that represents the root of the AST.
	Root NodeType = iota

	Fn
	Args
	Scope
	Type
	ReturnStmt
	Expression
)

// The AST node that generated from the parser.
type Node struct {
	typ   NodeType
	token *token.Token

	parent *Node
	childs []*Node
}

// Show the AST tree in the human-readable format as the tree.
func (n *Node) String() string {
	prefix, indent := "└── ", "    "
	return n.showIndentString(prefix, indent, 0)
}

func (n *Node) showIndentString(prefix, indent string, level int) string {
	builder := strings.Builder{}

	builder.WriteString(prefix)
	builder.WriteString(fmt.Sprintf("%s: %v", n.typ, n.token))

	level += 1

	for idx, child := range n.childs {
		builder.WriteString("\n")

		switch {
		case idx == len(n.childs)-1:
			pre := strings.Repeat(indent, level) + "└── "
			builder.WriteString(child.showIndentString(pre, indent, level))
		default:
			pre := strings.Repeat(indent, level) + "├── "
			builder.WriteString(child.showIndentString(pre, indent, level))
		}
	}

	return builder.String()
}

// Get the level of the AST node in the tree.
func (n *Node) Level() int {
	level := 0
	parent := n.parent

	for parent != nil {
		level++
		parent = parent.parent
	}

	return level
}

// Get the type of the AST node.
func (n *Node) Type() NodeType {
	return n.typ
}

// Get the raw token that generated the AST node.
func (n *Node) Token() *token.Token {
	return n.token
}

// Get the parent node of the AST node.
func (n *Node) Parent() *Node {
	return n.parent
}

// List the children nodes of the AST node.
func (n *Node) Children() []*Node {
	return n.childs
}

// Append a new child node into the tail of child to the AST node.
func (n *Node) Append(child *Node) {
	child.parent = n
	n.childs = append(n.childs, child)
}
