package lox

import (
	"fmt"
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

func (e *environment) assign(t token, value any) {
	_, ok := e.values[t.lexeme]
	if !ok {
		err := newError(fmt.Sprintf("Undefined variable %s.", t.lexeme), t.line)
		panic(err)
	}
	e.values[t.lexeme] = value
}

func (e *environment) get(t token) any {
	value, ok := e.values[t.lexeme]
	if !ok {
		err := newError(fmt.Sprintf("Undefined variable %s.", t.lexeme), t.line)
		panic(err)
	}
	return value
}
