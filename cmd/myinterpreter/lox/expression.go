package lox

import (
	"fmt"
)

type exprInterface interface {
	fmt.Stringer
	evaluate() string
	get() expression
}

type expression struct {
	expression exprInterface
	right      exprInterface
	operator   token
}

func (e *expression) String() string {
	if e == nil {
		return ""
	}
	if e.expression == nil && e.right == nil {
		return e.primary()
	}
	return e.parenthesized()
}

func (e *expression) evaluate() string {
	return ""
}

func (e *expression) get() expression {
	return *e
}

func (e *expression) primary() string {
	if e.operator.literal != NULL {
		return e.operator.literal
	}
	return e.operator.lexeme
}

func (e *expression) parenthesized() string {
	if e.expression == nil {
		return fmt.Sprintf("(%s %s)", e.operator.lexeme, e.right.String())
	}
	if e.right == nil {
		return fmt.Sprintf("(%s %s)", e.operator.lexeme, e.expression.String())
	}
	return fmt.Sprintf("(%s %s %s)", e.operator.lexeme, e.expression.String(), e.right.String())
}
