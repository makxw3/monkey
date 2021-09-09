package main

import (
	"fmt"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

func main() {
	input := `let a = 10 + 20;
	1 + 2 + 3 + 4 * 5 * 6 * 7;`
	lx := lexer.New(input)
	ps := parser.New(lx)

	program := ps.ParseProgram()
	if program != nil {
		for _, stmt := range program.Statements {
			fmt.Println(stmt.String())
		}
	}
}
