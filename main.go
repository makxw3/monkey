package main

import (
	"fmt"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

func main() {
	input := `fn(one,two,three,four,five) { 90; }`
	lx := lexer.New(input)
	ps := parser.New(lx)

	program := ps.ParseProgram()
	if program != nil {
		for _, stmt := range program.Statements {
			fmt.Printf("%v\n", stmt.String())
		}
	}
}
