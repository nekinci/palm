package main

type SyntaxKind int

const (
	BadToken SyntaxKind = iota
	EOFToken
	NumberToken
	PlusToken
	MinusToken
	StarToken
	SlashToken
	WhitespaceToken
)

func (sk SyntaxKind) String() string {
	switch sk {
	case BadToken:
		return "BadToken"
	case EOFToken:
		return "EOFToken"
	case NumberToken:
		return "NumberToken"
	case PlusToken:
		return "PlusToken"
	case MinusToken:
		return "MinusToken"
	case StarToken:
		return "StarToken"
	case SlashToken:
		return "SlashToken"
	case WhitespaceToken:
		return "WhitespaceToken"
	default:
		return string(sk)
	}
}

type Token struct {
	Text     string
	Position int
	Kind     SyntaxKind
	Value    interface{}
}

type Lexer struct {
	sourceCode string
	position   int
	start      int
	kind       SyntaxKind
}

func NewLexer(sourceCode string) *Lexer {
	return &Lexer{sourceCode: sourceCode}
}

func (l *Lexer) Do() Token {
	return l.readToken()
}

func (l *Lexer) readToken() Token {
	l.start = l.position
	l.kind = BadToken

	if l.position >= len(l.sourceCode) {
		l.kind = EOFToken
		return l.makeToken()
	}

	// Whitespace tokens
	if isWhitespace(l.Current()) {
		for isWhitespace(l.Current()) {
			l.position++
		}
		l.kind = WhitespaceToken
		return l.makeToken()
	}

	// Number tokens
	if isDigit(l.Current()) {
		for isDigit(l.Current()) {
			l.position++
		}
		l.kind = NumberToken
		return l.makeToken()
	}

	// Plus tokens
	if l.Current() == '+' {
		l.position++
		l.kind = PlusToken
		return l.makeToken()
	}

	// Minus tokens
	if l.Current() == '-' {
		l.position++
		l.kind = MinusToken
		return l.makeToken()
	}

	// Star tokens
	if l.Current() == '*' {
		l.position++
		l.kind = StarToken
		return l.makeToken()
	}

	// Slash tokens
	if l.Current() == '/' {
		l.position++
		l.kind = SlashToken
		return l.makeToken()
	}

	return l.makeToken()

}

func (l *Lexer) makeToken() Token {
	return Token{
		Text:     l.sourceCode[l.start:l.position],
		Position: l.start,
		Kind:     l.kind,
	}
}

func (l *Lexer) Peek(offset int) rune {
	if l.position+offset >= len(l.sourceCode) {
		return 0
	}
	return rune(l.sourceCode[l.position+offset])
}

func (l *Lexer) Current() rune {
	return l.Peek(0)
}

func (l *Lexer) Lookahead() rune {
	return l.Peek(1)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
