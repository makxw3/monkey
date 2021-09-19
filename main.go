package main

import (
	"fmt"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

func main() {
	input := `1 + 2 * 3 + 4`
	lx := lexer.New(input)
	ps := parser.New(lx)

	program := ps.ParseProgram()
	if program != nil {
		for _, stmt := range program.Statements {
			fmt.Printf("%v\n", stmt.String())
		}
	}
}

// Precedence
// The precedence value of an operator defines which operator should be evaluated first
