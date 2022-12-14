package parse

import "strconv"

type SyntaxTree struct {
	Root Node
}

type Parser struct {
	lexer     *Lexer
	pos       int
	token     Token
	lookAhead Token
	tree      *SyntaxTree
}

func NewParser(name, input string) *Parser {
	p := &Parser{
		lexer: NewLexer(name, input),
		tree:  &SyntaxTree{},
	}

	// ignore first token because it is nil
	_ = p.getCurrentAndNext()
	return p
}

// getCurrentAndNext advances the lexer next token and return current token
func (p *Parser) getCurrentAndNext() Token {
	currentToken := p.currentToken()
	p.token = p.lexer.nextToken()
	return currentToken
}

func (p *Parser) currentToken() Token {
	return p.token
}

func (p *Parser) Parse() (*SyntaxTree, error) {
	t := &SyntaxTree{}
	t.Root = p.parseExpression(0)
	return t, nil
}

func (p *Parser) parseExpression(parentPrecedence int) Node {
	var left Node
	unaryPrecedence := p.currentToken().Kind.GetUnaryPrecedence()

	if unaryPrecedence != 0 && unaryPrecedence >= parentPrecedence {
		token := p.getCurrentAndNext()
		expr := p.parseExpression(unaryPrecedence)
		left = NewUnaryExpressionNode(p.tree, token, expr)
	} else {
		left = p.parsePrimary()
	}

	for {
		precedence := p.currentToken().Kind.GetBinaryPrecedence()
		if precedence == 0 || precedence <= parentPrecedence {
			break
		}
		opToken := p.getCurrentAndNext()
		right := p.parseExpression(precedence)
		left = NewBinaryExpressionNode(p.tree, left, opToken, right)
	}

	return left
}

func (p *Parser) parsePrimary() Node {
	switch p.currentToken().Kind {
	case NUMBER:
		return p.parseNumber()
	case LPAREN:
		return p.parseParenthesizedExpression()
	case FALSE, TRUE:
		return p.parseBoolean()
	}

	return nil
}

func (p *Parser) parseParenthesizedExpression() Node {
	openParenthesisToken := p.getCurrentAndNext()
	expr := p.parseExpression(0)
	if p.currentToken().Kind != RPAREN {
		return nil
	}
	closeParenthesisToken := p.getCurrentAndNext()
	return NewParenthesizedExpressionNode(p.tree, openParenthesisToken, expr, closeParenthesisToken)
}

func (p *Parser) parseNumber() Node {
	val := p.getCurrentAndNext().Val
	// Todo think about error handling
	valInt, _ := strconv.Atoi(val)

	return &NumberNode{
		NodeKind:   NodeNumber,
		NumberKind: NumberInt,
		Raw:        val,
		Int:        int64(valInt),
	}
}

func (p *Parser) parseBoolean() Node {
	val := p.getCurrentAndNext().Val
	// Todo think about error handling
	valBool, _ := strconv.ParseBool(val)

	return &BooleanNode{
		NodeKind: NodeBoolean,
		Raw:      val,
		Val:      valBool,
	}
}
