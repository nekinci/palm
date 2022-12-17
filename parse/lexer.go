package parse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// This is Rob Pike style lexer based on the state functions and the state machine.
// I have adopted this style 'cause it's easier and more understandable and maintainable. So At least I think so.

type Location struct {
	Line  int
	Start int
	Len   int
	File  string
}

type StateFn func(*Lexer) StateFn

type TokenKind int

const endOfFile = -1

const (
	EOF TokenKind = iota // End of input
	ERROR
	PLUS         // +
	MINUS        // -
	MUL          // *
	QUO          // /
	REM          // %
	PLUS_ASSIGN  // +=
	MINUS_ASSIGN // -=
	MUL_ASSIGN   // *=
	QUO_ASSIGN   // /=
	REM_ASSIGN   // %=
	GT           // >
	LT           // <
	GTE          // >=
	LTE          // <=
	EQ           // ==
	NEQ          // !=
	ASSIGN       // =
	DECLARE      // :=
	COLON        // :
	BITAND       // &
	BITOR        // |
	AND          // &&
	OR           // ||
	XOR          // ^
	LSHIFT       // <<
	RSHIFT       // >>
	LPAREN       // (
	RPAREN       // )
	NUMBER       // 12345
	NOT          // !
	FALSE        // false
	TRUE         // true

)

var stringKind = map[string]TokenKind{
	"+":  PLUS,
	"-":  MINUS,
	"*":  MUL,
	"/":  QUO,
	"%":  REM,
	"(":  LPAREN,
	")":  RPAREN,
	"!":  NOT,
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	"<":  LT,
	">=": GTE,
	"<=": LTE,
	"=":  ASSIGN,
	":=": DECLARE,
	":":  COLON,
	"^":  XOR,
	"+=": PLUS_ASSIGN,
	"-=": MINUS_ASSIGN,
	"*=": MUL_ASSIGN,
	"/=": QUO_ASSIGN,
	"%=": REM_ASSIGN,
	"<<": LSHIFT,
	">>": RSHIFT,
}

func (k TokenKind) String() string {
	switch k {
	case EOF:
		return "EOF"
	case ERROR:
		return "ERROR"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MUL:
		return "MUL"
	case QUO:
		return "QUO"
	case REM:
		return "REM"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case NUMBER:
		return "NUMBER"
	default:
		return "UNKNOWN"
	}
}

func (k TokenKind) GetBinaryPrecedence() int {
	return GetBinaryOperatorPrecedence(k)
}

func (k TokenKind) GetUnaryPrecedence() int {
	return GetUnaryOperatorPrecedence(k)
}

type Token struct {
	Kind TokenKind
	Val  string
	pos  int
	len  int
	line int
}

func (t Token) Location(fileName string) Location {
	return Location{
		Line:  t.line,
		Start: t.pos,
		Len:   t.len,
		File:  fileName,
	}
}

type Lexer struct {
	name       string
	input      string
	start      int
	pos        int
	len        int
	line       int
	startLine  int
	tokens     chan Token
	diagnostic Diagnostic
}

func NewLexer(name, input string) *Lexer {
	l := &Lexer{
		name:   name,
		input:  input,
		tokens: make(chan Token),
	}
	go l.run()
	return l
}

func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// Next returns the next rune in the input
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.len = 0
		return endOfFile
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	// set length of new next rune that readen from input by old pos
	l.len = w
	l.pos += l.len

	if r == '\n' {
		l.line++
	}
	return r

}

// Backup recover to back step
func (l *Lexer) backup() {
	l.pos -= l.len
	if l.len == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// ignore ignores some input like whitespaces.
func (l *Lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	l.startLine = l.line

}

// accept advances rune if it's valid by given input.
// If it is valid then returns true or else returns false
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun advances runes until being invalid one.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {

	}
	l.backup()
}

// acceptExact advances rune if it's equal to given input.
func (l *Lexer) acceptExact(valid string) bool {
	if strings.HasPrefix(l.input[l.pos:], valid) {
		for range valid {
			l.next()
		}

		return true
	}
	return false
}

func (l *Lexer) errorf(format string, args ...any) StateFn {
	l.tokens <- Token{Kind: ERROR, Val: fmt.Sprintf(format, args), pos: l.start, len: l.line}
	return lexText
}

func (l *Lexer) emit(kind TokenKind) {
	l.tokens <- Token{
		Kind: kind,
		Val:  l.input[l.start:l.pos],
		pos:  l.pos,
		line: l.startLine,
		len:  l.pos - l.start,
	}

	l.start = l.pos
	l.startLine = l.line
}

func (l *Lexer) nextToken() Token {
	return <-l.tokens
}

// State Functions

func lexText(l *Lexer) StateFn {
	l.len = 0

	if l.pos >= len(l.input) {
		l.emit(EOF)
		return nil
	}

	switch r := l.peek(); {
	case r == endOfFile:
		l.emit(EOF)
		return nil
	case isWhitespace(r):
		return lexWhitespace
	case r == '(':
		return lexLeftParen
	case r == '+' || r == '-' || r == '*' || r == '/' || r == '%' || r == '!' || r == '<' ||
		r == '>' || r == '=' || r == ':' || r == '&' || r == '|' || r == '^':
		return lexOperator
	case r == ')':
		return lexRightParen
	case r >= '0' && r <= '9':
		return lexNumber
	// at least now we check only for t and f for true and false after that we will check for all keywords
	case r == 't' || r == 'f':
		return lexBoolean
	default:
		return l.errorf("unrecognized character in input: %q", r)
	}

}

func lexBoolean(l *Lexer) StateFn {
	if l.acceptExact("true") {
		l.emit(TRUE)
		return lexText
	}

	if l.acceptExact("false") {
		l.emit(FALSE)
		return lexText
	}

	return l.errorf("unrecognized character in input: %q", l.next())
}

func lexLeftParen(l *Lexer) StateFn {
	l.accept("(")
	l.emit(LPAREN)
	return lexText
}

func lexRightParen(l *Lexer) StateFn {
	l.accept(")")
	l.emit(RPAREN)
	return lexText
}

func lexNumber(l *Lexer) StateFn {
	l.acceptRun("0123456789")
	l.emit(NUMBER)
	return lexText
}

func lexOperator(l *Lexer) StateFn {

	if l.accept(">") {
		if l.accept("=") {
			l.emit(GTE)
			return lexText
		}

		if l.accept(">") {
			l.emit(RSHIFT)
			return lexText
		}

		l.emit(GT)
		return lexText
	}

	if l.accept("<") {
		if l.accept("=") {
			l.emit(LTE)
			return lexText
		}

		if l.accept("<") {
			l.emit(LSHIFT)
			return lexText
		}

		l.emit(LT)
		return lexText
	}

	if l.accept(":") {
		if l.accept("=") {
			l.emit(DECLARE)
			return lexText
		}

		// we don't handle yet, think about it later
		l.backup()
	}

	if l.accept("!") {
		if l.accept("=") {
			l.emit(NEQ)
			return lexText
		}

		l.emit(NOT)
		return lexText
	}

	if l.accept("&") {
		if l.accept("&") {
			l.emit(AND)
			return lexText
		}

		l.emit(BITAND)
	}

	if l.accept("|") {
		if l.accept("|") {
			l.emit(OR)
			return lexText
		}

		l.emit(BITOR)
	}

	if l.accept("^") {
		l.emit(XOR)
		return lexText
	}

	if l.accept("=") {
		if l.accept("=") {
			l.emit(EQ)
			return lexText
		}

		l.emit(ASSIGN)
		return lexText
	}

	if l.accept("+-*/%") {
		l.accept("=")
		l.emit(stringKind[l.input[l.start:l.pos]])
	}

	return lexText
}

func lexWhitespace(l *Lexer) StateFn {
	l.acceptRun(Whitespace)
	l.ignore()
	return lexText
}
