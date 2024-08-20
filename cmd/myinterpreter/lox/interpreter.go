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

func (i *interpreter) visitExpr(e *expression) string { return "" }

func (i *interpreter) visitEquality(e *expressionEquality) string {
	leftIsString := i.evaluatesToString(e.get().expression)
	rightIsString := i.evaluatesToString(e.get().right)
	bothString := leftIsString && rightIsString
	bothNonString := !leftIsString && !rightIsString
	haveSameType := bothString || bothNonString
	left := e.get().expression.accept(i)
	right := e.get().right.accept(i)
	switch e.get().operator.tokenType {
	case EQUAL_EQUAL:
		return fmt.Sprintf("%t", left == right && haveSameType)
	case BANG_EQUAL:
		return fmt.Sprintf("%t", left != right || !haveSameType)
	}
	return ""
}

func (i *interpreter) visitComparison(e *expressionComparison) string {
	left := i.getFloat(e.get().expression)
	right := i.getFloat(e.get().right)
	switch e.get().operator.tokenType {
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
	switch e.get().operator.tokenType {
	case PLUS:
		if i.evaluatesToString(e) {
			return e.get().expression.accept(i) + e.get().right.accept(i)
		}
		left := i.getFloat(e.get().expression)
		right := i.getFloat(e.get().right)
		return i.printFloat(left + right)
	case MINUS:
		left := i.getFloat(e.get().expression)
		right := i.getFloat(e.get().right)
		return i.printFloat(left - right)
	}
	return ""
}

func (i *interpreter) visitFactor(e *expressionFactor) string {
	left := i.getFloat(e.get().expression)
	right := i.getFloat(e.get().right)
	switch e.get().operator.tokenType {
	case STAR:
		return i.printFloat(left * right)
	case SLASH:
		return i.printFloat(left / right)
	}
	return ""
}

func (i *interpreter) visitUnary(e *expressionUnary) string {
	switch e.get().operator.tokenType {
	case BANG:
		return fmt.Sprintf("%t", !i.isTruthy(e.get().right))
	case MINUS:
		val := i.getFloat(e.get().right)
		return i.printFloat(-val)
	default:
		return ""
	}
}

func (i *interpreter) visitLiteral(e *expressionLiteral) string {
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

func (i *interpreter) visitGroup(e *expressionGroup) string {
	return e.exprInterface.accept(i)
}

func (i *interpreter) isTruthy(e exprInterface) bool {
	value := e.accept(i)
	switch e.get().operator.tokenType {
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
		i.runtimeErrors = append(i.runtimeErrors, newError("Operand must be a number.", e.get().operator.line))
	}
	return valDouble
}

func (i *interpreter) printFloat(val float64) string {
	reg1 := regexp.MustCompile(`([\d])0*$`)
	reg2 := regexp.MustCompile(`\.0*$`)
	return reg2.ReplaceAllString(reg1.ReplaceAllString(fmt.Sprintf("%f", val), "$1"), "")
}

func (i *interpreter) evaluatesToString(e exprInterface) bool {
	if e.get().expression == nil || e.get().right == nil {
		return e.get().operator.tokenType == STRING
	}
	if e.get().expression.get().operator.tokenType == PLUS &&
		e.get().right.get().operator.tokenType == PLUS {
		return i.evaluatesToString(e.get().expression) && i.evaluatesToString(e.get().right)
	}
	if e.get().expression.get().operator.tokenType == PLUS &&
		e.get().right.get().operator.tokenType == STRING {
		return i.evaluatesToString(e.get().expression)
	}
	if e.get().expression.get().operator.tokenType == STRING &&
		e.get().right.get().operator.tokenType == PLUS {
		return i.evaluatesToString(e.get().right)
	}
	return e.get().expression.get().operator.tokenType == STRING &&
		e.get().right.get().operator.tokenType == STRING
}
