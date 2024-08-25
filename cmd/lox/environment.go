package lox

import (
	"fmt"
	"os"
)

type environment struct {
	values map[string]any
}

func newEnvironment() *environment {
	return &environment{values: make(map[string]any)}
}

func (e *environment) define(name string, value any) {
	e.values[name] = value
}

func (e *environment) get(t token) any {
	value, ok := e.values[t.lexeme]
	if !ok {
		err := newError(fmt.Sprintf("Undefined variable %s.", t.lexeme), t.line)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(70)
	}
	return value
}
