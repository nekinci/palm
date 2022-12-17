package main

import (
	"fmt"
	"myProgrammingLanguage/parse"
)

func main() {
	parser := parse.NewParser("test", "1 << 1")

	evaluator := NewEvaluator(parser)
	result, _ := evaluator.Evaluate()
	fmt.Println(result)

}
