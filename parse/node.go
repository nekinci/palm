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
	NodeCallExpression
	NodeIfStatement
	NodeElseStatement
	NodeBlockStatement
	NodeVariableDeclaration
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
	TypeToken  Token
	Identifier Token
	Op         Token
	Right      Node
}

func NewAssignmentExpressionNode(tree *SyntaxTree, typeToken Token, identifier Token, op Token, right Node) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		NodeKind:   NodeAssignmentExpression,
		TypeToken:  typeToken,
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

type CallExpressionNode struct {
	NodeKind
	tr         *SyntaxTree
	Pos        int
	Identifier Token
}

func NewCallExpressionNode(tree *SyntaxTree, identifier Token) *CallExpressionNode {
	return &CallExpressionNode{
		NodeKind:   NodeCallExpression,
		Identifier: identifier,
		tr:         tree,
	}
}

func (n *CallExpressionNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *CallExpressionNode) String() string {
	return n.Identifier.Val
}

func (n *CallExpressionNode) Position() int {
	return n.Pos
}

func (n *CallExpressionNode) tree() *SyntaxTree {
	return n.tr
}

func (n *CallExpressionNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Identifier.Val)
}

///////////////////////////////////////////////////////////

type BlockStatementNode struct {
	NodeKind
	tr    *SyntaxTree
	Pos   int
	Left  Token
	Right Token
	Nodes []Node
}

func NewBlockStatementNode(tree *SyntaxTree, left Token, right Token, nodes []Node) *BlockStatementNode {
	return &BlockStatementNode{
		NodeKind: NodeBlockStatement,
		Left:     left,
		Right:    right,
		Nodes:    nodes,
		tr:       tree,
	}
}

func (n *BlockStatementNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *BlockStatementNode) String() string {
	// Think about this should we return with the nodes or not?
	return n.Left.Val + n.Right.Val
}

func (n *BlockStatementNode) Position() int {
	return n.Pos
}

func (n *BlockStatementNode) tree() *SyntaxTree {
	return n.tr
}

func (n *BlockStatementNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Left.Val)
	for _, node := range n.Nodes {
		builder.WriteString(node.String())
	}
	builder.WriteString(n.Right.Val)
}

///////////////////////////////////////////////////////////

type IfStatementNode struct {
	NodeKind
	tr         *SyntaxTree
	Pos        int
	IfToken    Token
	Expression Node
	Body       Node
	Else       Node
}

func NewIfStatementNode(tree *SyntaxTree, ifToken Token, expression Node, body Node, elseNode Node) *IfStatementNode {
	return &IfStatementNode{
		NodeKind:   NodeIfStatement,
		IfToken:    ifToken,
		Expression: expression,
		Body:       body,
		Else:       elseNode,
		tr:         tree,
	}
}

func (n *IfStatementNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *IfStatementNode) String() string {
	return n.IfToken.Val + n.Expression.String() + n.Body.String()
}

func (n *IfStatementNode) Position() int {
	return n.Pos
}

func (n *IfStatementNode) tree() *SyntaxTree {
	return n.tr
}

func (n *IfStatementNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.IfToken.Val)
	builder.WriteString(n.Expression.String())
	builder.WriteString(n.Body.String())
	if n.Else != nil {
		builder.WriteString(n.Else.String())
	}
}

///////////////////////////////////////////////////////////

type ElseStatementNode struct {
	NodeKind
	tr   *SyntaxTree
	Pos  int
	Else Token
	Body Node
}

func NewElseStatementNode(tree *SyntaxTree, elseToken Token, body Node) *ElseStatementNode {
	return &ElseStatementNode{
		NodeKind: NodeElseStatement,
		Else:     elseToken,
		Body:     body,
		tr:       tree,
	}
}

func (n *ElseStatementNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *ElseStatementNode) String() string {
	return n.Else.Val + n.Body.String()
}

func (n *ElseStatementNode) Position() int {
	return n.Pos
}

func (n *ElseStatementNode) tree() *SyntaxTree {
	return n.tr
}

func (n *ElseStatementNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.Else.Val)
	builder.WriteString(n.Body.String())
}

///////////////////////////////////////////////////////////

type VariableDeclarationStatementNode struct {
	NodeKind
	tr           *SyntaxTree
	Pos          int
	TypeToken    Token
	DeclareToken Token
	HasTypeToken bool
	Identifier   Token
	Expression   Node
}

func NewVariableDeclarationNode(tree *SyntaxTree, typeToken *Token, identifier Token, declareToken Token, expression Node) *VariableDeclarationStatementNode {

	hasTypeToken := false
	if typeToken != nil {
		hasTypeToken = true
	}

	node := &VariableDeclarationStatementNode{
		NodeKind:     NodeVariableDeclaration,
		HasTypeToken: hasTypeToken,
		Identifier:   identifier,
		Expression:   expression,
		DeclareToken: declareToken,
		tr:           tree,
	}

	if hasTypeToken {
		node.TypeToken = *typeToken
	}

	return node
}

func (n *VariableDeclarationStatementNode) Kind() NodeKind {
	return n.NodeKind
}

func (n *VariableDeclarationStatementNode) String() string {
	return n.TypeToken.Val + n.Identifier.Val + n.Expression.String()
}

func (n *VariableDeclarationStatementNode) Position() int {
	return n.Pos
}

func (n *VariableDeclarationStatementNode) tree() *SyntaxTree {
	return n.tr
}

func (n *VariableDeclarationStatementNode) writeTo(builder *strings.Builder) {
	builder.WriteString(n.TypeToken.Val)
	builder.WriteString(n.Identifier.Val)
	builder.WriteString(n.Expression.String())
}

///////////////////////////////////////////////////////////
