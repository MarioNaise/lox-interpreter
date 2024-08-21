package lox

import (
	"fmt"
	"regexp"
	"strconv"
)

type interpreter struct {
	runtimeErrors []loxError
}

func (i *interpreter) evaluate(expr exprInterface) (string, []loxError) {
	return expr.accept(i), i.runtimeErrors
}

func (i *interpreter) visitEquality(e *expressionEquality) string {
	leftIsString := i.evaluatesToString(e.expr())
	rightIsString := i.evaluatesToString(e.next())
	bothString := leftIsString && rightIsString
	bothNonString := !leftIsString && !rightIsString
	haveSameType := bothString || bothNonString
	left := e.expr().accept(i)
	right := e.next().accept(i)
	switch e.tokenType() {
	case EQUAL_EQUAL:
		return fmt.Sprintf("%t", left == right && haveSameType)
	case BANG_EQUAL:
		return fmt.Sprintf("%t", left != right || !haveSameType)
	}
	return ""
}

func (i *interpreter) visitComparison(e *expressionComparison) string {
	left := i.getFloat(e.expr())
	right := i.getFloat(e.next())
	switch e.tokenType() {
	case LESS:
		return fmt.Sprintf("%t", left < right)
	case LESS_EQUAL:
		return fmt.Sprintf("%t", left <= right)
	case GREATER:
		return fmt.Sprintf("%t", left > right)
	case GREATER_EQUAL:
		return fmt.Sprintf("%t", left >= right)
	}
	return ""
}

func (i *interpreter) visitTerm(e *expressionTerm) string {
	switch e.tokenType() {
	case PLUS:
		if i.evaluatesToString(e) {
			return e.expr().accept(i) + e.next().accept(i)
		}
		left := i.getFloat(e.expr())
		right := i.getFloat(e.next())
		return i.printFloat(left + right)
	case MINUS:
		left := i.getFloat(e.expr())
		right := i.getFloat(e.next())
		return i.printFloat(left - right)
	}
	return ""
}

func (i *interpreter) visitFactor(e *expressionFactor) string {
	left := i.getFloat(e.expr())
	right := i.getFloat(e.next())
	switch e.tokenType() {
	case STAR:
		return i.printFloat(left * right)
	case SLASH:
		return i.printFloat(left / right)
	}
	return ""
}

func (i *interpreter) visitUnary(e *expressionUnary) string {
	switch e.tokenType() {
	case BANG:
		return fmt.Sprintf("%t", !i.isTruthy(e.next()))
	case MINUS:
		val := i.getFloat(e.next())
		return i.printFloat(-val)
	default:
		return ""
	}
}

func (i *interpreter) visitLiteral(e *expressionLiteral) string {
	regexNr := regexp.MustCompile(`\.0$`)
	switch e.tokenType() {
	case STRING:
		return e.literal()
	case NUMBER:
		return regexNr.ReplaceAllString(e.literal(), "")
	default:
		return e.lexeme()
	}
}

func (i *interpreter) visitGroup(e *expressionGroup) string {
	return e.exprInterface.accept(i)
}

func (i *interpreter) visitExpr(e *expression) string { return "" }

func (i *interpreter) isTruthy(e exprInterface) bool {
	value := e.accept(i)
	switch e.tokenType() {
	case BANG:
		return value == "true"
	case NIL:
		return value != "nil"
	case NUMBER:
		return value != "0"
	case STRING:
		return value != ""
	case TRUE:
		return true
	}
	return false
}

func (i *interpreter) getFloat(e exprInterface) float64 {
	valDouble, err := strconv.ParseFloat(e.accept(i), 64)
	if err != nil {
		error := newError("Operand must be a number.", e.token().line)
		i.runtimeErrors = append(i.runtimeErrors, error)
	}
	return valDouble
}

func (i *interpreter) printFloat(val float64) string {
	asString := fmt.Sprintf("%f", val)
	// 2.230000 -> 2.23
	re := regexp.MustCompile(`([\d])0*$`)
	// 2.0 -> 2
	reg := regexp.MustCompile(`\.0*$`)
	str := re.ReplaceAllString(asString, "$1")
	return reg.ReplaceAllString(str, "")
}

func (i *interpreter) evaluatesToString(e exprInterface) bool {
	if e.expr() == nil || e.next() == nil {
		return e.tokenType() == STRING
	}
	if e.expr().tokenType() == PLUS &&
		e.next().tokenType() == PLUS {
		return i.evaluatesToString(e.expr()) && i.evaluatesToString(e.next())
	}
	if e.expr().tokenType() == PLUS &&
		e.next().tokenType() == STRING {
		return i.evaluatesToString(e.expr())
	}
	if e.expr().tokenType() == STRING &&
		e.next().tokenType() == PLUS {
		return i.evaluatesToString(e.next())
	}
	return e.expr().tokenType() == STRING &&
		e.next().tokenType() == STRING
}
