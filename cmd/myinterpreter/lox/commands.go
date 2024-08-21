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
	expr, parseErrors := p.parse()
	if len(parseErrors) == 0 {
		printer := astPrinter{}
		printer.print(expr)
	}
	printErrors(parseErrors)
	return len(parseErrors) == 0
}

func Evaluate(r io.Reader) bool {
	p := newParser(r)
	expr, parseErrors := p.parse()
	if len(parseErrors) == 0 {
		i := interpreter{}
		result, errors := i.evaluate(expr)
		if len(errors) == 0 {
			fmt.Println(result)
		} else {
			printErrors(errors)
			return false
		}
	}
	printErrors(parseErrors)
	return len(parseErrors) == 0
}

func printErrors(errors []loxError) {
	for _, error := range errors {
		fmt.Fprintln(os.Stderr, error)
	}
}
