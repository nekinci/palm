package main

import "myProgrammingLanguage/parse"

type Interpreter struct {
	scope *parse.Scope
	tree  *parse.SyntaxTree
}
