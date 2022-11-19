package main

type Expression struct {
	token Token
}

type PlusExpression struct {
	Expression
	left  Expression
	right Expression
}

type Parser struct {
	lexer             *Lexer
	tokens            []Token
	primaryExpression *Expression
	currentToken      int
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	for {
		token := p.lexer.Do()
		if token.Kind == EOFToken {
			break
		}
		p.tokens = append(p.tokens, token)
	}

	return p
}

func (p *Parser) Parse() {
	// Parse the tokens
	p.primaryExpression = p.ParseExpression()
}

func (p *Parser) ParseExpression() *Expression {
	// Parse the tokens
	return nil
}
