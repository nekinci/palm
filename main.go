package main

import (
	"bufio"
	"fmt"
	"myProgrammingLanguage/parse"
	"os"
)

func main() {
	fmt.Println("Welcome to my programming language!")
	repl()
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	text := ""
	scope := parse.NewScope(nil)
	for {
		fmt.Print(">> ")
		text, _ = reader.ReadString('\n')
		if text == "exit\n" {
			fmt.Println("Bye!")
			break
		}
		parser := parse.NewParser("repl", text)
		tree, _ := parser.Parse()
		evaluator := NewEvaluator(tree, scope)
		result, _ := evaluator.Evaluate()
		fmt.Println(result)
	}
}
