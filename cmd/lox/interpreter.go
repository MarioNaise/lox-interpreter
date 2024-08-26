package lox

import (
	"fmt"
	"io"
	"reflect"
)

type interpreter struct {
	*parser
	*environment
}

func newInterpreter(r io.Reader) *interpreter {
	env := newEnvironment()
	p := newParser(r)
	return &interpreter{p, env}
}

func (i *interpreter) evaluate(e exprInterface) any {
	return e.accept(i)
}

func (i *interpreter) execute(s stmtInterface) {
	s.accept(i)
}

func (i *interpreter) interpret(stmts []stmtInterface) {
	for _, s := range stmts {
		i.handleStmt(s)
	}
}

func (i *interpreter) handleStmt(s stmtInterface) {
	defer func() {
		if r := recover(); r != nil {
			i.synchronize()
			panic(r)
		}
	}()
	s.accept(i)
}

func (i *interpreter) visitVarStmt(s *stmtVar) {
	var val any
	var name string
	if s.expr() != nil {
		val = i.evaluate(s.expr())
		name = s.name().lexeme
	}
	i.define(name, val)
}

func (i *interpreter) visitPrintStmt(s *stmtPrint) {
	val := i.evaluate(s.expr())
	fmt.Printf("%v\n", val)
}

func (i *interpreter) visitExprStmt(s *stmtExpr) {
	i.evaluate(s.expr())
}

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
			return fmt.Sprintf("%v%v", e.expr().accept(i), e.next().accept(i))
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
		return !i.isTruthy(e.next())
	case MINUS:
		val := i.parseFloat(e.next())
		return -val
	default:
		return false
	}
}

func (i *interpreter) visitVar(e *expressionVar) any {
	return i.get(e.token())
}

func (i *interpreter) visitAssignment(e *expressionAssignment) any {
	i.assign(e.expr().token(), i.evaluate(e.next()))
	return i.get(e.expr().token())
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
	left := reflect.TypeOf(e.expr().accept(i)).Name()
	right := reflect.TypeOf(e.next().accept(i)).Name()
	return left == "string" ||
		right == "string"
}

func (i *interpreter) parseFloat(e exprInterface) float64 {
	n := i.evaluate(e)
	if reflect.TypeOf(n).Name() != "float64" {
		err := newError(fmt.Sprintf("Operand must be a number: %v", e.lexeme()), e.token().line)
		panic(err)
	}
	return n.(float64)
}

func (i *interpreter) hasSameType(a any, b any) bool {
	return reflect.TypeOf(a).Name() == reflect.TypeOf(b).Name()
}

func (i *interpreter) isTruthy(e exprInterface) bool {
	value := e.value()
	if value == nil {
		return false
	}
	switch e.tokenType() {
	case IDENTIFIER:
		expr := &expression{val: i.get(e.token())}
		return i.parser.isTruthy(expr)
	case BANG:
		return e.accept(i).(bool)
	case STRING:
		return value != ""
	case NUMBER:
		return value.(float64) != 0
	case TRUE:
		return true
	case FALSE:
		return false
	case NIL:
		return false
	}
	return false
}
