package lox

import (
	"fmt"
	"io"
	"os"
)

func Tokenize(r io.Reader) bool {
	s := newScanner(r)
	s.tokenize()
	for _, token := range s.tokens {
		fmt.Println(token)
	}
	printErrors(s.scanErrors)
	return len(s.scanErrors) == 0
}

func Parse(r io.Reader) bool {
	p := newParser(r)
	p.parse()
	if len(p.parseErrors) == 0 {
		printer := astPrinter{}
		printer.print(p.expression)
	}
	printErrors(p.parseErrors)
	return len(p.parseErrors) == 0 && len(p.scanErrors) == 0
}

func Evaluate(r io.Reader) bool {
	p := newParser(r)
	p.parse()
	if len(p.parseErrors) == 0 {
		i := interpreter{}
		result, errors := i.evaluate(p.expression)
		if len(errors) == 0 {
			fmt.Println(result)
		} else {
			printErrors(i.runtimeErrors)
			return false
		}
	}
	printErrors(p.scanErrors)
	printErrors(p.parseErrors)
	return len(p.parseErrors) == 0 && len(p.scanErrors) == 0
}

func printErrors(errors []loxError) {
	for _, error := range errors {
		fmt.Fprintln(os.Stderr, error)
	}
}
