package main

import "fmt"

func main() {

	var lexer *Lexer = NewLexer("(1 + 2)")
	var token Token = lexer.Do()

	for token.Kind != EOFToken {
		if token.Kind != WhitespaceToken {
			fmt.Printf("TokenKind: %v, TokenText: %v\n", token.Kind.String(), token.Text)
		}
		token = lexer.Do()
	}
}
