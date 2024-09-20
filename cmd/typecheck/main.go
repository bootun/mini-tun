package main

import (
	"fmt"

	"github.com/bootun/mini-tun/pkg/lexer"
	"github.com/bootun/mini-tun/pkg/parser"
	"github.com/bootun/mini-tun/pkg/typecheck"
)

func main() {
	input := `let a = 1
let b = 20
let c = a + b
let add = function(a, b) {
	let c = a + b - 1
	return c + 10 - b
}
let e = add(10,add(a))
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

	if err := typecheck.NewChecker(program).Check(); err != nil {
		fmt.Printf("type check error: %v\n", err)
		return
	}
	fmt.Printf("pass type check")

}
