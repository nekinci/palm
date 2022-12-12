package main

import (
	"fmt"
	"myProgrammingLanguage/parse"
)

func main() {
	parser := parse.NewParser("test", "(1+2) * (3+1) + 5")

	evaluator := NewEvaluator(parser)
	result, _ := evaluator.Evaluate()
	fmt.Println(result)
}
