package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// This is Rob Pike style lexer based on the state functions and the state machine.
// I have adopted this style 'cause it's easier and more understandable and maintainable. So At least I think so.

type StateFn func(*Lexer) StateFn

type TokenKind int

const endOfFile = -1

const (
	EOF TokenKind = iota // End of input
	BADTOKEN
	UNEXPECTED   // Unexpected token
	EMPTY        // Empty token means that there is no token
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
	IDENT        // main
	INT          // int
	IF           // if
	ELSE         // else
	LBRACE       // {
	RBRACE       // }
	BOOL         // bool
)

var emptyToken = Token{
	Kind: EMPTY,
	Val:  "",
	len:  0,
}

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
	case BADTOKEN:
		return "BADTOKEN"
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
	case NOT:
		return "NOT"
	case EQ:
		return "EQ"
	case NEQ:
		return "NEQ"
	case GT:
		return "GT"
	case LT:
		return "LT"
	case GTE:
		return "GTE"
	case LTE:
		return "LTE"
	case ASSIGN:
		return "ASSIGN"
	case DECLARE:
		return "DECLARE"
	case COLON:
		return "COLON"
	case XOR:
		return "XOR"
	case PLUS_ASSIGN:
		return "PLUS_ASSIGN"
	case MINUS_ASSIGN:
		return "MINUS_ASSIGN"
	case MUL_ASSIGN:
		return "MUL_ASSIGN"
	case QUO_ASSIGN:
		return "QUO_ASSIGN"
	case REM_ASSIGN:
		return "REM_ASSIGN"
	case LSHIFT:
		return "LSHIFT"
	case RSHIFT:
		return "RSHIFT"
	case IDENT:
		return "IDENT"
	case FALSE:
		return "FALSE"
	case TRUE:
		return "TRUE"
	case INT:
		return "INT"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case BOOL:
		return "BOOL"
	default:
		panic(fmt.Sprintf("unknown token kind: %d", k))
	}
}

func (k TokenKind) GetBinaryPrecedence() int {
	return GetBinaryOperatorPrecedence(k)
}

func (k TokenKind) GetUnaryPrecedence() int {
	return GetUnaryOperatorPrecedence(k)
}

type TokenLocation struct {
	Start Location
	End   Location
}

type Location struct {
	Offset   int
	Line     int
	Col      int
	Filename string
}

func (l Location) String() string {
	return fmt.Sprintf("Location(Filename = %s, Line = %d, Column = %d, Offset = %d)", l.Filename, l.Line, l.Col, l.Offset)
}

type Token struct {
	Kind TokenKind
	Val  string
	len  int
	Loc  TokenLocation
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.Kind, t.Val)
}

type Lexer struct {
	name          string
	input         string
	startOffset   int
	offset        int
	startLineCols map[int]int
	lineCols      map[int]int
	len           int
	line          int
	startLine     int
	tokens        chan Token
	errors        chan Err
	doneErr       chan bool
}

func NewLexer(name, input string) *Lexer {
	l := &Lexer{
		name:          name,
		input:         input,
		tokens:        make(chan Token),
		lineCols:      make(map[int]int),
		startLineCols: make(map[int]int),
		errors:        make(chan Err),
		doneErr:       make(chan bool),
	}
	go l.run()
	return l
}

func (l *Lexer) loc() TokenLocation {
	return TokenLocation{
		Start: Location{
			Offset:   l.startOffset,
			Line:     l.startLine,
			Filename: l.name,
			Col:      l.startLineCols[l.startLine],
		},
		End: Location{
			Offset:   l.offset,
			Line:     l.line,
			Filename: l.name,
			Col:      l.lineCols[l.line],
		},
	}
}

func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	l.doneErr <- true
	close(l.errors)
	close(l.tokens)
}

// Next returns the next rune in the input
func (l *Lexer) next() rune {
	if l.offset >= len(l.input) {
		l.len = 0
		return endOfFile
	}
	r, w := utf8.DecodeRuneInString(l.input[l.offset:])
	// set length of new next rune that readen from input by old offset
	l.len = w
	l.offset += l.len
	l.lineCols[l.line] += l.len
	if r == '\n' {
		l.line++
	}
	return r

}

// Backup recover to back step
func (l *Lexer) backup() {
	l.offset -= l.len
	if l.len == 1 && l.input[l.offset] == '\n' {
		l.line--
	} else {
		l.lineCols[l.line] -= l.len
	}
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// ignore ignores some input like whitespaces.
func (l *Lexer) ignore() {
	// l.line += strings.Count(l.input[l.startOffset:l.offset], "\n")
	l.startOffset = l.offset
	l.startLineCols[l.line] = l.lineCols[l.line]
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

func (l *Lexer) acceptRunFunc(f func(rune) bool) {
	for f(l.next()) {

	}
	l.backup()
}

// acceptExact advances rune if it's equal to given input.
func (l *Lexer) acceptExact(valid string) bool {
	if strings.HasPrefix(l.input[l.offset:], valid) {
		for range valid {
			l.next()
		}

		return true
	}
	return false
}

func (l *Lexer) errorf(kind TokenKind, format string, args ...any) StateFn {
	l.emit(kind)
	l.errors <- Err{
		Kind: Error,
		Len:  l.len,
		Msg:  fmt.Sprintf(format, args...),
		File: l.name,
		Loc:  l.loc(),
	}
	l.next()
	return lexText
}

func (l *Lexer) emit(kind TokenKind) {
	l.tokens <- Token{
		Kind: kind,
		Val:  l.input[l.startOffset:l.offset],
		len:  l.offset - l.startOffset,
		Loc:  l.loc(),
	}

	l.startOffset = l.offset
	l.startLine = l.line
	l.startLineCols[l.line] = l.lineCols[l.line]
}

func (l *Lexer) nextToken() Token {
	return <-l.tokens
}

// State Functions

func lexText(l *Lexer) StateFn {
	l.len = 0

	if l.offset >= len(l.input) {
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
	case r == '{':
		return lexLeftBrace
	case r == '}':
		return lexRightBrace
	default:
		if unicode.IsLetter(r) {
			return lexIdentifierOrKeyword
		}
		return l.errorf(BADTOKEN, "unrecognized character in input: %q", r)
	}

}

func lexLeftBrace(l *Lexer) StateFn {
	l.next()
	l.emit(LBRACE)
	return lexText
}

func lexRightBrace(l *Lexer) StateFn {
	l.next()
	l.emit(RBRACE)
	return lexText
}

func lexIdentifierOrKeyword(l *Lexer) StateFn {
	l.acceptRunFunc(unicode.IsLetter)
	tok := l.input[l.startOffset:l.offset]

	if tok == "true" {
		l.emit(TRUE)
		return lexText
	}

	if tok == "false" {
		l.emit(FALSE)
		return lexText
	}

	if tok == "if" {
		l.emit(IF)
		return lexText
	}

	if tok == "else" {
		l.emit(ELSE)
		return lexText
	}

	if tok == "int" {
		l.emit(INT)
		return lexText
	}

	if tok == "bool" {
		l.emit(BOOL)
		return lexText
	}

	l.emit(IDENT)
	return lexText
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
		l.emit(stringKind[l.input[l.startOffset:l.offset]])
	}

	return lexText
}

func lexWhitespace(l *Lexer) StateFn {
	l.acceptRun(Whitespace)
	l.ignore()
	return lexText
}
