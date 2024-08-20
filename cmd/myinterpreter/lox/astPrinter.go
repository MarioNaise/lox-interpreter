package lox

import (
	"fmt"
	"strings"
)

type astPrinter struct{}

func (a *astPrinter) print(expr exprInterface) {
	fmt.Println(expr.accept(a))
}

func (a *astPrinter) visitExpr(e *expression) string { return "" }

func (a *astPrinter) visitEquality(e *expressionEquality) string {
	return a.parenthesized(e.get().operator.lexeme, e.get().expression, e.get().right)
}

func (a *astPrinter) visitComparison(e *expressionComparison) string {
	return a.parenthesized(e.get().operator.lexeme, e.get().expression, e.get().right)
}

func (a *astPrinter) visitTerm(e *expressionTerm) string {
	return a.parenthesized(e.get().operator.lexeme, e.get().expression, e.get().right)
}

func (a *astPrinter) visitFactor(e *expressionFactor) string {
	return a.parenthesized(e.get().operator.lexeme, e.get().expression, e.get().right)
}

func (a *astPrinter) visitUnary(e *expressionUnary) string {
	return a.parenthesized(e.get().operator.lexeme, e.get().right)
}

func (a *astPrinter) visitLiteral(e *expressionLiteral) string {
	return a.primary(e.get())
}

func (a *astPrinter) visitGroup(e *expressionGroup) string {
	return a.parenthesized("group", e.exprInterface)
}

func (a *astPrinter) primary(e *expression) string {
	if e.operator.literal != NULL {
		return e.operator.literal
	}
	return e.operator.lexeme
}

func (a *astPrinter) parenthesized(name string, e ...exprInterface) string {
	s := []string{}
	s = append(s, name)
	for _, expr := range e {
		if expr != nil {
			s = append(s, expr.accept(a))
		}
	}
	str := strings.Join(s, " ")
	return strings.Join([]string{"(", str, ")"}, "")
}
