package main

import (
	"fmt"

	"github.com/bootun/mini-tun/pkg/lexer"
	"github.com/bootun/mini-tun/pkg/parser"
)

func main() {
	input := `let a = 1
let b = 2
let c = a + b
let add = function(a, b) {
	return a + b
}
let e = add(1,add(a,b))
`
	l := lexer.New(input)
	newParser, err := parser.New(l)
	if err != nil {
		panic(err)
	}
	program, err := newParser.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Printf("program AST: %v", program.JSON())
	// for _, statement := range program.Statements {
	//	fmt.Println(statement.TokenLiteral())
	// }
}
