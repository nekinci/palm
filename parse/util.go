package parse

const Whitespace = " \t\r\n"

func GetUnaryOperatorPrecedence(kind TokenKind) int {
	switch kind {
	case PLUS, MINUS, NOT:
		return 6
	}
	return 0
}

func GetBinaryOperatorPrecedence(kind TokenKind) int {
	switch kind {
	case PLUS, MINUS:
		return 4
	case MUL, QUO, REM, LSHIFT, RSHIFT:
		return 5
	case GT, LT, GTE, LTE, NEQ, EQ:
		return 3
	case AND, BITAND:
		return 2
	case OR, BITOR, XOR:
		return 1
	}
	return 0
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func getMapValueOrDefineAndGetValue(m map[int]int, key int) interface{} {
	if _, ok := m[key]; !ok {
		m[key] = 0
	}
	return m[key]
}
