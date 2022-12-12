package parse

import "strings"

type NodeKind int

type NumberKind int

type Node interface {
	Kind() NodeKind
	String() string
	Position() int
	tree() *SyntaxTree
	writeTo(builder *strings.Builder)
}

const (
	NodeEOF NodeKind = iota
	NodeBinaryExpression
	NodeNumber
	NodeParenthesisedExpression
)

const (
	NumberInt NumberKind = iota
	NumberFloat
	NumberComplex
	NumberUnsigned
)

///////////////////////////////////////////////////////////

// NumberNode TODO implement float, complex and unsigned integers
type NumberNode struct {
	NodeKind
	NumberKind
	tr  *SyntaxTree
	Pos int
	Raw string
	Int int64
}

func (n *NumberNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *NumberNode) String() string {
	return n.Raw
}

func (n *NumberNode) Position() int {
	return n.Pos
}

func (n *NumberNode) tree() *SyntaxTree {
	return n.tr
}

func (n *NumberNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Raw)
}

///////////////////////////////////////////////////////////

type BinaryExpressionNode struct {
	NodeKind
	tr    *SyntaxTree
	Pos   int
	Left  Node
	Op    Token
	Right Node
}

func NewBinaryExpressionNode(tree *SyntaxTree, left Node, op Token, right Node) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		NodeKind: NodeBinaryExpression,
		Left:     left,
		Op:       op,
		Right:    right,
		tr:       tree,
	}
}

func (n *BinaryExpressionNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *BinaryExpressionNode) String() string {
	return n.Left.String() + n.Op.Val + n.Right.String()
}

func (n *BinaryExpressionNode) Position() int {
	return n.Pos
}

func (n *BinaryExpressionNode) tree() *SyntaxTree {
	return n.tr
}

func (n *BinaryExpressionNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Left.String())
	builder.WriteString(n.Op.Val)
	builder.WriteString(n.Right.String())
}

///////////////////////////////////////////////////////////

type ParenthesisedExpressionNode struct {
	NodeKind
	tr         *SyntaxTree
	Left       Token
	Expression Node
	Right      Token
	Pos        int
}

func NewParenthesizedExpressionNode(tree *SyntaxTree, left Token, expression Node, right Token) *ParenthesisedExpressionNode {
	return &ParenthesisedExpressionNode{
		NodeKind:   NodeParenthesisedExpression,
		Left:       left,
		Expression: expression,
		Right:      right,
		tr:         tree,
	}
}

func (n *ParenthesisedExpressionNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *ParenthesisedExpressionNode) String() string {
	return n.Left.Val + n.Expression.String() + n.Right.Val
}

func (n *ParenthesisedExpressionNode) Position() int {
	return n.Pos
}

func (n *ParenthesisedExpressionNode) tree() *SyntaxTree {
	return n.tr
}

func (n *ParenthesisedExpressionNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Left.Val)
	builder.WriteString(n.Expression.String())
	builder.WriteString(n.Right.Val)
}

///////////////////////////////////////////////////////////
