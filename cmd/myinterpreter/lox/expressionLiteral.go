package lox

import (
	"regexp"
)

type expressionLiteral struct {
	exprInterface
}

func (e *expressionLiteral) evaluate() string {
	regexNr := regexp.MustCompile(`\.0$`)
	operator := e.get().operator
	switch operator.tokenType {
	case STRING:
		return operator.literal
	case NUMBER:
		return regexNr.ReplaceAllString(operator.literal, "")
	default:
		return operator.lexeme
	}
}
