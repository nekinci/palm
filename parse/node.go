package parse

import (
	"strconv"
	"strings"
)

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
	NodeBoolean
	NodeParenthesisedExpression
	NodeUnaryExpression
	NodeAssignmentExpression
	NodeIdentifierAccessExpression
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

type BooleanNode struct {
	NodeKind
	tr  *SyntaxTree
	Pos int
	Raw string
	Val bool
}

func NewBooleanNode(tr *SyntaxTree, val bool, pos int) *BooleanNode {
	return &BooleanNode{
		NodeKind: NodeNumber,
		tr:       tr,
		Pos:      pos,
		Raw:      strconv.FormatBool(val),
		Val:      val,
	}
}

func (n *BooleanNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *BooleanNode) String() string {
	return n.Raw
}

func (n *BooleanNode) Position() int {
	return n.Pos
}

func (n *BooleanNode) tree() *SyntaxTree {
	return n.tr
}

func (n *BooleanNode) writeTo(builder *strings.Builder) {
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

type UnaryExpressionNode struct {
	NodeKind
	tr    *SyntaxTree
	Pos   int
	Op    Token
	Right Node
}

func NewUnaryExpressionNode(tree *SyntaxTree, op Token, right Node) *UnaryExpressionNode {
	return &UnaryExpressionNode{
		NodeKind: NodeUnaryExpression,
		Op:       op,
		Right:    right,
		tr:       tree,
	}
}

func (n *UnaryExpressionNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *UnaryExpressionNode) String() string {
	return n.Op.Val + n.Right.String()
}

func (n *UnaryExpressionNode) Position() int {
	return n.Pos
}

func (n *UnaryExpressionNode) tree() *SyntaxTree {
	return n.tr
}

func (n *UnaryExpressionNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Op.Val)
	builder.WriteString(n.Right.String())
}

///////////////////////////////////////////////////////////

type AssignmentExpressionNode struct {
	NodeKind
	tr         *SyntaxTree
	Pos        int
	Identifier Token
	Op         Token
	Right      Node
}

func NewAssignmentExpressionNode(tree *SyntaxTree, identifier Token, op Token, right Node) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		NodeKind:   NodeAssignmentExpression,
		Identifier: identifier,
		Op:         op,
		Right:      right,
		tr:         tree,
	}
}

func (n *AssignmentExpressionNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *AssignmentExpressionNode) String() string {
	return n.Identifier.Val + n.Op.Val + n.Right.String()
}

func (n *AssignmentExpressionNode) Position() int {
	return n.Pos
}

func (n *AssignmentExpressionNode) tree() *SyntaxTree {
	return n.tr
}

func (n *AssignmentExpressionNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Identifier.Val)
	builder.WriteString(n.Op.Val)
	builder.WriteString(n.Right.String())
}

///////////////////////////////////////////////////////////

type IdentifierAccessExpressionNode struct {
	NodeKind
	tr         *SyntaxTree
	Pos        int
	Identifier Token
}

func NewIdentifierAccessExpressionNode(tree *SyntaxTree, identifier Token) *IdentifierAccessExpressionNode {
	return &IdentifierAccessExpressionNode{
		NodeKind:   NodeIdentifierAccessExpression,
		Identifier: identifier,
		tr:         tree,
	}
}

func (n *IdentifierAccessExpressionNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *IdentifierAccessExpressionNode) String() string {
	return n.Identifier.Val
}

func (n *IdentifierAccessExpressionNode) Position() int {
	return n.Pos
}

func (n *IdentifierAccessExpressionNode) tree() *SyntaxTree {
	return n.tr
}

func (n *IdentifierAccessExpressionNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Identifier.Val)
}

///////////////////////////////////////////////////////////
