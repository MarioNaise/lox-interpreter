package lox

import (
	"fmt"
	"os"
)

// turns a string into a list of tokens and errors, returns true if no errors
func Tokenize(str string) bool {
	s := newLoxScanner()
	s.tokenize(str)
	printTokens(s.tokens)
	printErrors(s.errors)
	return len(s.errors) == 0
}

// tokenizes a string, then parses the tokens into an expression
// returns true if no errors
func Parse(str string) bool {
	p := parser{}
	p.parse(str)
	if len(p.parseErrors) == 0 {
		fmt.Println(p.expression)
	}
	for _, error := range p.parseErrors {
		fmt.Fprintln(os.Stderr, error)
	}
	return len(p.parseErrors) == 0 && len(p.errors) == 0
}

// prints a list of tokens
func printTokens(tokens []token) {
	for _, token := range tokens {
		fmt.Println(token)
	}
}

// prints a list of errors
func printErrors(errors []loxError) {
	for _, error := range errors {
		fmt.Fprintln(os.Stderr, error)
	}
}
