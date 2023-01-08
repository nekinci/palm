package parse

import (
	"strconv"
	"sync"
)

type SyntaxTree struct {
	Root Node
}

type Parser struct {
	lexer     *Lexer
	pos       int
	tokens    []Token
	badTokens []Token
	tree      *SyntaxTree
	Errors    ErrorContainer
}

func NewParser(name, input string) *Parser {
	p := &Parser{
		lexer: NewLexer(name, input),
		tree:  &SyntaxTree{},
		Errors: ErrorContainer{
			Errors: []Err{},
			mu:     &sync.Mutex{},
		},
		badTokens: []Token{},
	}

	go func() {
		// consume errors from lexer until it closes
		for {
			select {
			case err := <-p.lexer.errors:
				p.Errors.AddError(err)
			case done := <-p.lexer.doneErr:
				if done {
					return
				}
			}
		}
	}()

	return p
}

// getCurrentAndNext advances the lexer next token and return current token
func (p *Parser) getCurrentAndNext() Token {

	token := p.currentToken()
	p.pos++
	return token

}

func (p *Parser) expect2(kinds ...TokenKind) Token {
	// iterate all kinds
	for _, kind := range kinds {
		if p.currentToken().Kind == kind {
			return p.getCurrentAndNext()
		}
	}

	kindsStr := ""
	for _, kind := range kinds {
		kindsStr += kind.String() + ", "
	}

	p.Errors.AddError(Err{
		File: p.lexer.name,
		Len:  p.currentToken().len,
		Loc:  p.currentToken().Loc,
		Msg:  "expected one of them <" + kindsStr + "> got <" + p.currentToken().Kind.String() + ">",
		Kind: Error,
	})

	return Token{
		Kind: UNEXPECTED,
		len:  p.currentToken().len,
		Val:  p.currentToken().Val,
		Loc:  p.currentToken().Loc,
	}
}

func (p *Parser) expect(kind TokenKind) Token {
	if p.currentToken().Kind == kind {
		return p.getCurrentAndNext()
	}

	p.Errors.AddError(Err{
		File: p.lexer.name,
		Len:  p.currentToken().len,
		Loc:  p.currentToken().Loc,
		Msg:  "expected " + kind.String() + " got " + p.currentToken().Kind.String(),
		Kind: Error,
	})

	return Token{
		Kind: UNEXPECTED,
		len:  p.currentToken().len,
		Val:  p.currentToken().Val,
		Loc:  p.currentToken().Loc,
	}

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
	t.Root = p.parseStatement()
	return t, nil
}

func (p *Parser) parseStatement() Node {
	switch p.currentToken().Kind {
	case IF:
		return p.parseIfStatement()
	case LBRACE:
		return p.parseBlockStatement()
	case IDENT:
		if p.peek(1).Kind == DECLARE {
			return p.parseVariableDeclaration()
		}
		return p.parseExpression()
	case INT:
		return p.parseVariableDeclaration()
	case BADTOKEN:
		p.badTokens = append(p.badTokens, p.getCurrentAndNext())
		return p.parseStatement()
	}

	return p.parseExpression()
}

func (p *Parser) parseIfStatement() Node {
	ifToken := p.expect(IF)
	condition := p.parseExpression()
	body := p.parseStatement()
	elseBody := p.parseElseStatement()

	if elseBody != nil {

	}

	return NewIfStatementNode(p.tree, ifToken, condition, body, elseBody)

}

func (p *Parser) parseElseStatement() Node {
	token := p.getCurrentAndNext()
	if token.Kind != ELSE {
		return nil
	}

	// else if cases
	if p.currentToken().Kind == IF {
		return p.parseIfStatement()
	}

	elseBody := p.parseStatement()
	return NewElseStatementNode(p.tree, token, elseBody)
}

func (p *Parser) parseBlockStatement() Node {
	token := p.expect(LBRACE)
	statements := []Node{}
	for p.currentToken().Kind != RBRACE {
		statements = append(statements, p.parseStatement())
	}
	return NewBlockStatementNode(p.tree, token, p.expect(RBRACE), statements)
}

func (p *Parser) parseExpression() Node {
	if p.currentToken().Kind == IDENT {
		switch p.peek(1).Kind {
		case ASSIGN, PLUS_ASSIGN, MINUS_ASSIGN, MUL_ASSIGN, QUO_ASSIGN, REM_ASSIGN:
			return p.parseAssignmentExpression()
		}
	}

	return p.parseBinaryExpression(0)
}

func (p *Parser) parseVariableDeclaration() Node {
	if p.currentToken().Kind == IDENT && p.peek(1).Kind == DECLARE {
		ident := p.expect(IDENT)
		declareToken := p.expect(DECLARE)
		expr := p.parseExpression()
		return NewVariableDeclarationNode(p.tree, nil, ident, declareToken, expr)
	}

	typeToken := p.expect2(INT, BOOL)
	ident := p.expect(IDENT)
	declareToken := p.expect(ASSIGN)
	expr := p.parseExpression()
	return NewVariableDeclarationNode(p.tree, &typeToken, ident, declareToken, expr)

}

func (p *Parser) parseAssignmentExpression() Node {
	left := p.expect(IDENT)
	opToken := p.expect2(ASSIGN, PLUS_ASSIGN, MINUS_ASSIGN, MUL_ASSIGN, QUO_ASSIGN, REM_ASSIGN)
	right := p.parseExpression()

	return NewAssignmentExpressionNode(p.tree, emptyToken, left, opToken, right)

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
	ident := p.expect(IDENT)
	if p.currentToken().Kind == LPAREN {
		// return p.parseCall(ident)
	}

	return NewCallExpressionNode(p.tree, ident)
}

func (p *Parser) parseParenthesizedExpression() Node {
	openParenthesisToken := p.expect(LPAREN)
	expr := p.parseBinaryExpression(0)
	if p.currentToken().Kind != RPAREN {
		return nil
	}
	closeParenthesisToken := p.expect(RPAREN)
	return NewParenthesizedExpressionNode(p.tree, openParenthesisToken, expr, closeParenthesisToken)
}

func (p *Parser) parseNumber() Node {
	val := p.expect(NUMBER)
	valInt, err := strconv.Atoi(val.Val)

	if err != nil {
		p.Errors.AddError(Err{
			File: p.lexer.name,
			Len:  val.len,
			Loc:  val.Loc,
			Msg:  "Unable to parse number: " + val.Val,
			Kind: Error,
		})

		return &NumberNode{
			NodeKind:   NodeNumber,
			NumberKind: NumberInt,
			Raw:        val.Val,
		}
	}

	return &NumberNode{
		NodeKind:   NodeNumber,
		NumberKind: NumberInt,
		Raw:        val.Val,
		Int:        int64(valInt),
	}
}

func (p *Parser) parseBoolean() Node {
	val := p.expect2(FALSE, TRUE)
	valBool, err := strconv.ParseBool(val.Val)

	if err != nil {
		p.Errors.AddError(Err{
			File: p.lexer.name,
			Len:  val.len,
			Loc:  val.Loc,
			Msg:  "Unable to parse boolean: " + val.Val,
			Kind: Error,
		})

		return &BooleanNode{
			NodeKind: NodeBoolean,
			Raw:      val.Val,
		}
	}

	return &BooleanNode{
		NodeKind: NodeBoolean,
		Raw:      val.Val,
		Val:      valBool,
	}
}
