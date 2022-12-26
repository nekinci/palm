package parse

import "strconv"

type SyntaxTree struct {
	Root Node
}

type Parser struct {
	lexer  *Lexer
	pos    int
	tokens []Token
	tree   *SyntaxTree
}

func NewParser(name, input string) *Parser {
	p := &Parser{
		lexer: NewLexer(name, input),
		tree:  &SyntaxTree{},
	}

	return p
}

// getCurrentAndNext advances the lexer next token and return current token
func (p *Parser) getCurrentAndNext() Token {

	token := p.currentToken()
	p.pos++
	return token

}

func (p *Parser) currentToken() Token {

	if p.pos >= len(p.tokens) {
		p.tokens = append(p.tokens, <-p.lexer.tokens)
	}

	return p.tokens[p.pos]
}

func (p *Parser) peek(offset int) Token {
	for len(p.tokens) <= p.pos+offset {
		p.tokens = append(p.tokens, <-p.lexer.tokens)
	}
	return p.tokens[p.pos+offset]
}

func (p *Parser) Parse() (*SyntaxTree, error) {
	t := &SyntaxTree{}
	t.Root = p.parseExpression()
	return t, nil
}

func (p *Parser) parseExpression() Node {
	if p.currentToken().Kind == IDENT {
		switch p.peek(1).Kind {
		case ASSIGN:
			return p.parseIdentifier()
		}
	}

	return p.parseBinaryExpression(0)
}

func (p *Parser) parseIdentifier() Node {
	left := p.getCurrentAndNext()
	opToken := p.getCurrentAndNext()
	right := p.parseExpression()

	return NewAssignmentExpressionNode(p.tree, left, opToken, right)

}

func (p *Parser) parseBinaryExpression(parentPrecedence int) Node {
	var left Node
	unaryPrecedence := p.currentToken().Kind.GetUnaryPrecedence()

	if unaryPrecedence != 0 && unaryPrecedence >= parentPrecedence {
		token := p.getCurrentAndNext()
		expr := p.parseBinaryExpression(unaryPrecedence)
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
		right := p.parseBinaryExpression(precedence)
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
	case IDENT:
		return p.parseIdentifierAccessOrCall()
	}

	return nil
}

func (p *Parser) parseIdentifierAccessOrCall() Node {
	ident := p.getCurrentAndNext()
	if p.currentToken().Kind == LPAREN {
		// return p.parseCall(ident)
	}

	return NewIdentifierAccessExpressionNode(p.tree, ident)
}

func (p *Parser) parseParenthesizedExpression() Node {
	openParenthesisToken := p.getCurrentAndNext()
	expr := p.parseBinaryExpression(0)
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
