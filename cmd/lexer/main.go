package main

import (
	"fmt"

	"github.com/bootun/mini-tun/pkg/lexer"
)

func main() {
	input := `let a = 1
let b = 20
let add = function(a, b) {
	return a + b
}
let e = add(a,b)`
	lex := lexer.New(input)
	tokens, err := lex.Parse()
	if err != nil {
		panic(err)
	}
	for _, token := range tokens {
		fmt.Printf("{Type: %s, Literal: %s}\n", token.GetType(), token.GetLiteral())
	}
}
