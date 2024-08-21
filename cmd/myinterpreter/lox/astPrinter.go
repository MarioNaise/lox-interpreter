package lox

import (
	"fmt"
	"strings"
)

type astPrinter struct{}

func (a *astPrinter) print(expr exprInterface) {
	if expr == nil {
		return
	}
	fmt.Println(expr.accept(a))
}

func (a *astPrinter) visitEquality(e *expressionEquality) string {
	return a.defaultString(e)
}

func (a *astPrinter) visitComparison(e *expressionComparison) string {
	return a.defaultString(e)
}

func (a *astPrinter) visitTerm(e *expressionTerm) string {
	return a.defaultString(e)
}

func (a *astPrinter) visitFactor(e *expressionFactor) string {
	return a.defaultString(e)
}

func (a *astPrinter) visitUnary(e *expressionUnary) string {
	return a.parenthesized(e.lexeme(), e.next())
}

func (a *astPrinter) visitLiteral(e *expressionLiteral) string {
	return a.primary(e)
}

func (a *astPrinter) visitGroup(e *expressionGroup) string {
	return a.parenthesized("group", e.exprInterface)
}

func (a *astPrinter) visitExpr(e *expression) string { return "" }

func (a *astPrinter) primary(e exprInterface) string {
	if e.literal() != NULL {
		return e.literal()
	}
	return e.lexeme()
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

func (a *astPrinter) defaultString(e exprInterface) string {
	return a.parenthesized(e.lexeme(), e.expr(), e.next())
}
