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
		printer := astPrinter{}
		for _, s := range stmts {
			printer.print(s)
		}
	}
	printErrors(parseErrors)
	return len(parseErrors) == 0
}

func Evaluate(r io.Reader) bool {
	p := newParser(r)
	stmts, parseErrors := p.parse()
	if len(parseErrors) == 0 {
		return evaluateStatements(stmts)
	}
	printErrors(parseErrors)
	return len(parseErrors) == 0
}

func evaluateStatements(stmts []stmtInterface) bool {
	i := interpreter{}
	for _, stmt := range stmts {
		result, errors := i.evaluate(stmt)
		if len(errors) == 0 {
			fmt.Println(result)
		} else {
			printErrors(errors)
			return false
		}
	}
	return len(i.runtimeErrors) == 0
}

func printErrors(errors []loxError) {
	for _, error := range errors {
		fmt.Fprintln(os.Stderr, error)
	}
}
