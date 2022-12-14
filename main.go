package main

import (
	"fmt"
	"myProgrammingLanguage/parse"
)

func main() {
	parser := parse.NewParser("test", "!true")

	evaluator := NewEvaluator(parser)
	result, _ := evaluator.Evaluate()
	fmt.Println(result)
}
