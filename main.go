package main

import (
	"bufio"
	"fmt"
	"myProgrammingLanguage/parse"
	"os"
	"strings"
)

func main() {
	evalFile()
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

		trimmedText := strings.TrimSpace(text)
		if strings.HasPrefix(trimmedText, "if") {
			fmt.Println("if statements are not supported yet")
			continue
		}

		parser := parse.NewParser("repl", text)
		tree, _ := parser.Parse()
		evaluator := NewEvaluator(tree, scope)
		result, err := evaluator.Evaluate()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
}

func evalFile() {
	// Read from file
	file, err := os.Open("test.pd")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	// Read until EOF
	text := ""
	for {
		line, err := reader.ReadString('\n')
		text += line
		if err != nil {
			break
		}
	}

	// Parse
	parser := parse.NewParser("test.pd", text)
	tree, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	if parser.Errors.HasErrors() {
		parser.Errors.Print()
		return
	}

	// Evaluate
	evaluator := NewEvaluator(tree, parse.NewScope(nil))

	result, err := evaluator.Evaluate()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

}
