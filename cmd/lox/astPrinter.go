package lox

import (
	"fmt"
	"strings"
)

type astPrinter struct{}

func (a *astPrinter) print(e exprInterface) {
	if e == nil {
		return
	}
	fmt.Println(e.accept(a))
}

func (a *astPrinter) visitPrintStmt(s *stmtPrint) {
	a.print(s.expr())
}

func (a *astPrinter) visitExprStmt(s *stmtExpr) {
	a.print(s.expr())
}

func (a *astPrinter) visitEquality(e *expressionEquality) any {
	return a.defaultString(e)
}

func (a *astPrinter) visitComparison(e *expressionComparison) any {
	return a.defaultString(e)
}

func (a *astPrinter) visitTerm(e *expressionTerm) any {
	return a.defaultString(e)
}

func (a *astPrinter) visitFactor(e *expressionFactor) any {
	return a.defaultString(e)
}

func (a *astPrinter) visitUnary(e *expressionUnary) any {
	return a.parenthesized(e.lexeme(), e.next())
}

func (a *astPrinter) visitLiteral(e *expressionLiteral) any {
	return a.primary(e)
}

func (a *astPrinter) visitGroup(e *expressionGroup) any {
	return a.parenthesized("group", e.expr())
}

func (a *astPrinter) visitExpr(e *expression) any { return "" }

func (a *astPrinter) primary(e exprInterface) any {
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
			s = append(s, fmt.Sprintf("%v", expr.accept(a)))
		}
	}
	str := strings.Join(s, " ")
	return strings.Join([]string{"(", str, ")"}, "")
}

func (a *astPrinter) defaultString(e exprInterface) string {
	return a.parenthesized(e.lexeme(), e.expr(), e.next())
}
