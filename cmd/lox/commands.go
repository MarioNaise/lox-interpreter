package lox

import (
	"fmt"
	"io"
	"os"
)

func Tokenize(r io.Reader) bool {
	s := newScanner(r)
	tokens, errs := s.tokenize()
	for _, token := range tokens {
		fmt.Println(token)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Parse(r io.Reader) bool {
	p := newParser(r)
	stmts, parseErrors := p.parse()
	if len(parseErrors) == 0 {
		p := astPrinter{}
		p.print(stmts)
	}
	printErrors(parseErrors)
	return len(parseErrors) == 0
}

func Evaluate(r io.Reader) bool {
	i := newInterpreter(r)
	stmts, parseErrors := i.parse()
	if len(parseErrors) == 0 {
		i.interpret(stmts)
		return true
	}
	printErrors(parseErrors)
	return false
}

func printErrors(errors []loxError) {
	for _, error := range errors {
		fmt.Fprintln(os.Stderr, error)
	}
}
