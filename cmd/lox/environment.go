package lox

import (
	"fmt"
)

type environment struct {
	enclosing *environment
	values    map[string]any
}

func newEnvironment(env *environment) *environment {
	return &environment{env, make(map[string]any)}
}

func (e *environment) define(name string, value any) {
	e.values[name] = value
}

func (e *environment) assign(t token, value any) {
	_, ok := e.values[t.lexeme]
	if ok {
		e.values[t.lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(t, value)
		return
	}
	err := newError(fmt.Sprintf("Undefined variable %s.", t.lexeme), t.line)
	panic(err)
}

func (e *environment) get(t token) any {
	value, ok := e.values[t.lexeme]
	if ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.get(t)
	}
	err := newError(fmt.Sprintf("Undefined variable %s.", t.lexeme), t.line)
	panic(err)
}
