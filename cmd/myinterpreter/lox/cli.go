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
	for _, error := range s.scanErrors {
		fmt.Fprintln(os.Stderr, error)
	}
	return len(s.scanErrors) == 0
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
	return len(p.parseErrors) == 0 && len(p.scanErrors) == 0
}

func Evaluate(r io.Reader) bool {
	p := newParser(r)
	p.parse()
	if len(p.parseErrors) == 0 {
		fmt.Println(p.expression.evaluate())
	}
	for _, error := range p.scanErrors {
		fmt.Fprintln(os.Stderr, error)
	}
	for _, error := range p.parseErrors {
		fmt.Fprintln(os.Stderr, error)
	}
	return len(p.parseErrors) == 0 && len(p.scanErrors) == 0
}
