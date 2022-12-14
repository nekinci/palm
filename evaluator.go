package main

import (
	"myProgrammingLanguage/parse"
)

type Evaluator struct {
	p    *parse.Parser
	tree *parse.SyntaxTree
}

func NewEvaluator(p *parse.Parser) *Evaluator {
	e := Evaluator{p: p}
	e.tree, _ = p.Parse()
	return &e
}

func (e *Evaluator) Evaluate() (interface{}, error) {
	return e.visitNode(e.tree.Root)
}

func (e *Evaluator) visitNode(node parse.Node) (interface{}, error) {
	switch node.Kind() {
	case parse.NodeNumber:
		return e.visitNumberNode(node.(*parse.NumberNode)), nil
	case parse.NodeBoolean:
		return e.visitBooleanNode(node.(*parse.BooleanNode)), nil
	case parse.NodeBinaryExpression:
		return e.visitBinaryExpressionNode(node.(*parse.BinaryExpressionNode))
	case parse.NodeParenthesisedExpression:
		return e.visitNode(node.(*parse.ParenthesisedExpressionNode).Expression)
	case parse.NodeUnaryExpression:
		return e.visitUnaryExpressionNode(node.(*parse.UnaryExpressionNode))
	}
	return nil, nil
}

func (e *Evaluator) visitNumberNode(node *parse.NumberNode) interface{} {
	return node.Int
}

func (e *Evaluator) visitBooleanNode(node *parse.BooleanNode) bool {
	return node.Val
}

func (e *Evaluator) visitBinaryExpressionNode(node *parse.BinaryExpressionNode) (interface{}, error) {
	left, err := e.visitNode(node.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.visitNode(node.Right)
	if err != nil {
		return nil, err
	}

	switch node.Op.Kind {
	case parse.PLUS:
		return left.(int64) + right.(int64), nil
	case parse.MINUS:
		return left.(int64) - right.(int64), nil
	case parse.MUL:
		return left.(int64) * right.(int64), nil
	case parse.QUO:
		return left.(int64) / right.(int64), nil
	}

	return nil, nil
}

func (e *Evaluator) visitUnaryExpressionNode(node *parse.UnaryExpressionNode) (interface{}, error) {
	right, err := e.visitNode(node.Right)
	if err != nil {
		return nil, err
	}

	switch node.Op.Kind {
	case parse.PLUS:
		return right.(int64), nil
	case parse.MINUS:
		return -right.(int64), nil
	case parse.BANG:
		val, err := e.visitNode(node.Right)
		return !val.(bool), err
	}

	return nil, nil
}