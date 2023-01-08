package main

import (
	"fmt"
	"myProgrammingLanguage/parse"
)

type Evaluator struct {
	tree  *parse.SyntaxTree
	scope *parse.Scope
}

func (e *Evaluator) popScope() {
	e.scope = e.scope.Parent()
}

func (e *Evaluator) pushScope() {
	e.scope = parse.NewScope(e.scope)
}

// peekScope looks up the scope chain for the given parentStep if it is not fount returns top level
func (e *Evaluator) peekScope(parentStep int) *parse.Scope {
	scope := e.scope
	for i := 0; i < parentStep; i++ {
		if scope.Parent() == nil {
			scope = scope.Parent()
		} else {
			break
		}
	}
	return scope
}

func NewEvaluator(tree *parse.SyntaxTree, scope *parse.Scope) *Evaluator {
	e := Evaluator{tree: tree, scope: scope}
	return &e
}

func (e *Evaluator) Evaluate() (interface{}, error) {
	if e.tree == nil || e.tree.Root == nil {
		return nil, nil
	}

	return e.visitNode(e.tree.Root)
}

func (e *Evaluator) visitNode(node parse.Node) (interface{}, error) {
	switch node.Kind() {
	case parse.NodeBlockStatement:
		return e.visitBlockStatementNode(node.(*parse.BlockStatementNode))
	case parse.NodeIfStatement:
		return e.visitIfStatementNode(node.(*parse.IfStatementNode))
	case parse.NodeElseStatement:
		return e.visitElseStatementNode(node.(*parse.ElseStatementNode))
	case parse.NodeVariableDeclaration:
		return e.visitVariableDeclarationNode(node.(*parse.VariableDeclarationStatementNode))
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
	case parse.NodeAssignmentExpression:
		return e.visitAssignmentExpressionNode(node.(*parse.AssignmentExpressionNode))
	case parse.NodeCallExpression:
		return e.visitIdentifierAccessExpressionNode(node.(*parse.CallExpressionNode))
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
	case parse.REM:
		return left.(int64) % right.(int64), nil
	case parse.AND:
		// if they are boolean then convert bool and compare
		// or else convert int
		leftBool, leftIsBool := left.(bool)
		rightBool, rightIsBool := right.(bool)
		if leftIsBool && rightIsBool {
			return leftBool && rightBool, nil
		}
		return left.(int64) & right.(int64), nil
	case parse.OR:
		leftBool, leftIsBool := left.(bool)
		rightBool, rightIsBool := right.(bool)
		if leftIsBool && rightIsBool {
			return leftBool || rightBool, nil
		}
		return left.(int64) | right.(int64), nil
	case parse.EQ:
		return left == right, nil
	case parse.NEQ:
		return left != right, nil
	case parse.LT:
		return left.(int64) < right.(int64), nil
	case parse.LTE:
		return left.(int64) <= right.(int64), nil
	case parse.GT:
		return left.(int64) > right.(int64), nil
	case parse.GTE:
		return left.(int64) >= right.(int64), nil
	case parse.BITOR:
		return left.(int64) | right.(int64), nil
	case parse.BITAND:
		return left.(int64) & right.(int64), nil
	case parse.XOR:
		return left.(int64) ^ right.(int64), nil
	case parse.LSHIFT:
		return left.(int64) << right.(int64), nil
	case parse.RSHIFT:
		return left.(int64) >> right.(int64), nil

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
	case parse.NOT:
		val, err := e.visitNode(node.Right)
		return !val.(bool), err
	}

	return nil, nil
}

func (e *Evaluator) visitAssignmentExpressionNode(node *parse.AssignmentExpressionNode) (interface{}, error) {
	resolvedVal, ok := e.scope.Resolve(node.Identifier.Val)

	switch node.Op.Kind {
	case parse.ASSIGN:

		// for repl, it is no necessary to check if the variable is declared
		//if !ok {
		//	return nil, fmt.Errorf("identifier %s not found", node.Identifier.Val)
		//}

		val, err := e.visitNode(node.Right)
		e.scope.Define(node.Identifier.Val, val)
		return val, err
	case parse.PLUS_ASSIGN, parse.MINUS_ASSIGN, parse.MUL_ASSIGN, parse.QUO_ASSIGN, parse.REM_ASSIGN:
		if !ok {
			return nil, fmt.Errorf("variable %s not defined", node.Identifier.Val)
		}
		val, err := e.visitNode(node.Right)
		if err != nil {
			return nil, err
		}

		var result int64
		if resolvedValInt, ok := resolvedVal.(int64); ok {
			if resolvedVal != nil {
				switch node.Op.Kind {
				case parse.PLUS_ASSIGN:
					result = resolvedValInt + val.(int64)
				case parse.MINUS_ASSIGN:
					result = resolvedValInt - val.(int64)
				case parse.MUL_ASSIGN:
					result = resolvedValInt * val.(int64)
				case parse.QUO_ASSIGN:
					result = resolvedValInt / val.(int64)
				case parse.REM_ASSIGN:
					result = resolvedValInt % val.(int64)
				}
			}
		} else {
			return nil, fmt.Errorf("variable %s is not an integer", node.Identifier.Val)
		}

		e.scope.Define(node.Identifier.Val, result)
		return result, nil

	}
	return nil, nil
}

func (e *Evaluator) visitIdentifierAccessExpressionNode(node *parse.CallExpressionNode) (interface{}, error) {
	val, ok := e.scope.Resolve(node.Identifier.Val)
	if !ok {
		return nil, fmt.Errorf("%d:%d undefined variable %s", node.Identifier.Val)
	}

	return val, nil
}

//////

func (e *Evaluator) visitIfStatementNode(node *parse.IfStatementNode) (interface{}, error) {
	condition, err := e.visitNode(node.Expression)
	if err != nil {
		return nil, err
	}

	if condition.(bool) {
		return e.visitNode(node.Body)
	} else if node.Else != nil {
		return e.visitNode(node.Else)
	}

	return nil, nil
}

func (e *Evaluator) visitElseStatementNode(node *parse.ElseStatementNode) (interface{}, error) {
	return e.visitNode(node.Body)
}

func (e *Evaluator) visitBlockStatementNode(node *parse.BlockStatementNode) (interface{}, error) {
	e.pushScope()
	var response any
	for _, statement := range node.Nodes {
		val, err := e.visitNode(statement)
		if err != nil {
			return nil, err
		}
		response = val
	}

	e.popScope()
	return response, nil
}

func (e *Evaluator) visitVariableDeclarationNode(node *parse.VariableDeclarationStatementNode) (interface{}, error) {
	var val interface{}
	var err error
	if node.Expression != nil {
		val, err = e.visitNode(node.Expression)
		if err != nil {
			return nil, err
		}
	}

	if _, ok := e.scope.ResolveLocal(node.Identifier.Val); ok {
		return nil, fmt.Errorf("variable %s already defined", node.Identifier.Val)
	}

	if node.HasTypeToken {
		// if the type is specified, check if the value is of that type
		if node.TypeToken.Val == "int" {
			if _, ok := val.(int64); !ok {
				return nil, fmt.Errorf("variable %s is not an integer", node.Identifier.Val)
			}
		}

		if node.TypeToken.Val == "bool" {
			if _, ok := val.(bool); !ok {
				return nil, fmt.Errorf("variable %s is not a boolean", node.Identifier.Val)
			}
		}
	}

	e.scope.Define(node.Identifier.Val, val)
	return val, nil
}
