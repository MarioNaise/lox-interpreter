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

// prints a list of tokens
func printTokens(tokens []token) {
	for _, token := range tokens {
		fmt.Println(token.ToString())
	}
}

// prints a list of errors
func printErrors(errors []loxError) {
	for _, error := range errors {
		fmt.Fprintln(os.Stderr, error.ToString())
	}
}
