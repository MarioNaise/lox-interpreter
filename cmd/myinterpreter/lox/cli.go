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
	for _, error := range s.errors {
		fmt.Fprintln(os.Stderr, error)
	}
	return len(s.errors) == 0
}

func Parse(r io.Reader) bool {
	p := newParser(r)
	p.parse()
	if len(p.parseErrors) == 0 {
		fmt.Println(p.expression)
	}
	for _, error := range p.parseErrors {
		fmt.Fprintln(os.Stderr, error)
	}
	return len(p.parseErrors) == 0 && len(p.errors) == 0
}
