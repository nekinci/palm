package parse

const Whitespace = " \t\r\n"

func GetBinaryOperatorPrecedence(kind TokenKind) int {
	switch kind {
	case PLUS, MINUS:
		return 4
	case MUL, QUO:
		return 5
	}
	return 0
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
