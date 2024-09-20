package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/bootun/mini-tun/pkg/interpreter"
	"github.com/bootun/mini-tun/pkg/lexer"
	"github.com/bootun/mini-tun/pkg/parser"
	"github.com/bootun/mini-tun/pkg/typecheck"
)

var (
	help = flag.Bool("help", false, "show help")
)

func init() {
	flag.Parse()
}

func main() {
	if *help {
		flag.Usage()
		return
	} else if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	fileName := os.Args[1]
	sourceCode, err := os.Open(fileName)
	if err != nil {
		log.Printf("failed to open %v: %v", fileName, err)
		os.Exit(1)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(sourceCode); err != nil {
		log.Printf("failed to read source code: %v", err)
		os.Exit(2)
	}
	input := buf.String()
	l := lexer.New(input)
	newParser, err := parser.New(l)
	if err != nil {
		log.Printf("failed to create parser: %v", err)
		os.Exit(3)
	}
	program, err := newParser.Parse()
	if err != nil {
		log.Printf("failed to parse program: %v", err)
		os.Exit(4)
	}

	if err := typecheck.NewChecker(program).Check(); err != nil {
		log.Printf("type check error: %v\n", err)
		os.Exit(5)
	}
	// fmt.Printf("pass type check")

	vm := interpreter.NewInterpreter(program)
	if err := vm.Exec(); err != nil {
		log.Printf("execute error: %v\n", err)
		os.Exit(6)
	}

}
