package lox

import "fmt"

type loxError struct {
	message string
	line    int
}

func newError(message string, line int) loxError {
	return loxError{message: message, line: line}
}

func (e loxError) ToString() string {
	return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}
