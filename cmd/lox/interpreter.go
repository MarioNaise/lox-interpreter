package lox

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type interpreter struct{}

func (i *interpreter) evaluate(e exprInterface) any {
	return e.accept(i)
}

func (i *interpreter) execute(s stmtInterface) {
	s.accept(i)
}

func (i *interpreter) interpret(stmts []stmtInterface) {
	for _, s := range stmts {
		s.accept(i)
	}
}

func (i *interpreter) visitPrintStmt(s *stmtPrint) {
	val := i.evaluate(s.expr())
	fmt.Printf("%v\n", val)
}

func (i *interpreter) visitExprStmt(s *stmtExpr) {}

func (i *interpreter) visitEquality(e *expressionEquality) any {
	left := e.expr().accept(i)
	right := e.next().accept(i)
	if !i.hasSameType(left, right) {
		return false
	}
	switch e.tokenType() {
	case EQUAL_EQUAL:
		return left == right
	case BANG_EQUAL:
		return left != right
	}
	return e.value()
}

func (i *interpreter) visitComparison(e *expressionComparison) any {
	left := i.parseFloat(e.expr())
	right := i.parseFloat(e.next())
	switch e.tokenType() {
	case LESS:
		return left < right
	case LESS_EQUAL:
		return left <= right
	case GREATER:
		return left > right
	case GREATER_EQUAL:
		return left >= right
	}
	return e.value()
}

func (i *interpreter) visitTerm(e *expressionTerm) any {
	switch e.tokenType() {
	case PLUS:
		if i.evaluatesToString(e) {
			return fmt.Sprintf("%v%v", e.expr().value(), e.next().value())
		}
		left := i.parseFloat(e.expr())
		right := i.parseFloat(e.next())
		return left + right
	case MINUS:
		left := i.parseFloat(e.expr())
		right := i.parseFloat(e.next())
		return left - right
	}
	return 0
}

func (i *interpreter) visitFactor(e *expressionFactor) any {
	left := i.parseFloat(e.expr())
	right := i.parseFloat(e.next())
	switch e.tokenType() {
	case STAR:
		return left * right
	case SLASH:
		return left / right
	}
	return ""
}

func (i *interpreter) visitUnary(e *expressionUnary) any {
	switch e.tokenType() {
	case BANG:
		return !e.value().(bool)
	case MINUS:
		val := e.next().value().(float64)
		return -val
	default:
		return false
	}
}

func (i *interpreter) visitLiteral(e *expressionLiteral) any {
	return e.value()
}

func (i *interpreter) visitGroup(e *expressionGroup) any {
	return e.exprInterface.accept(i)
}

func (i *interpreter) visitExpr(e *expression) any {
	return ""
}

func (i *interpreter) evaluatesToString(e exprInterface) bool {
	if e.expr() == nil || e.next() == nil {
		return e.tokenType() == STRING
	}
	if e.expr().tokenType() == PLUS &&
		e.next().tokenType() == PLUS {
		return i.evaluatesToString(e.expr()) && i.evaluatesToString(e.next())
	}
	return e.expr().tokenType() == STRING ||
		e.next().tokenType() == STRING
}

func (i *interpreter) parseFloat(e exprInterface) float64 {
	asFloat, err := strconv.ParseFloat(fmt.Sprintf("%v", e.value()), 64)
	if err != nil || reflect.TypeOf(e.value()).Name() != "float64" {
		error := newError(fmt.Sprintf("Operand must be a number: %v", e.lexeme()), e.token().line)
		fmt.Fprintln(os.Stderr, error)
		os.Exit(70)
	}
	return asFloat
}

func (i *interpreter) hasSameType(a any, b any) bool {
	return reflect.TypeOf(a).Name() == reflect.TypeOf(b).Name()
}
