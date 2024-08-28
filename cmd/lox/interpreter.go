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
	env := newEnvironment(nil)
	p := newParser(r)
	return &interpreter{p, env}
}

func (i *interpreter) evaluate(e exprInterface) any {
	return e.accept(i)
}

func (i *interpreter) execute(s stmtInterface) {
	defer func() {
		if r := recover(); r != nil {
			i.synchronize()
			panic(r)
		}
	}()
	s.accept(i)
}

func (i *interpreter) interpret(stmts []stmtInterface) {
	for _, s := range stmts {
		i.execute(s)
	}
}

func (i *interpreter) visitVarStmt(s *stmtVar) {
	var val any
	name := s.name().lexeme
	if s.expr() != nil {
		val = i.evaluate(s.expr())
	}
	i.define(name, val)
}

func (i *interpreter) visitIfStmt(s *stmtIf) {
	if i.isTruthy(s.condition) {
		i.execute(s.thenBranch)
	} else if s.elseBranch != nil {
		i.execute(s.elseBranch)
	}
}

func (i *interpreter) visitPrintStmt(s *stmtPrint) {
	val := i.evaluate(s.expr())
	fmt.Printf(i.stringify(val) + "\n")
}

func (i *interpreter) visitBlockStmt(s *stmtBlock) {
	prevEnv := i.environment
	i.environment = newEnvironment(prevEnv)
	i.interpret(s.statements)
	i.environment = prevEnv
}

func (i *interpreter) visitExprStmt(s *stmtExpr) {
	i.evaluate(s.expr())
}

func (i *interpreter) visitVar(e *expressionVar) any {
	return i.get(e.token())
}

func (i *interpreter) visitAssignment(e *expressionAssignment) any {
	i.assign(e.expr().token(), i.evaluate(e.next()))
	return i.get(e.expr().token())
}

func (i *interpreter) visitLogical(e *expressionLogical) any {
	if e.tokenType() == OR {
		if i.isTruthy(e.expr()) {
			return i.evaluate(e.expr())
		}
	} else if !i.isTruthy(e.expr()) {
		return i.evaluate(e.expr())
	}
	return i.evaluate(e.next())
}

func (i *interpreter) visitEquality(e *expressionEquality) any {
	left := i.evaluate(e.expr())
	right := i.evaluate(e.next())
	if !i.hasSameType(left, right) {
		return false
	}
	switch e.tokenType() {
	case EQUAL_EQUAL:
		return left == right
	case BANG_EQUAL:
		return left != right
	}
	return nil
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
	return nil
}

func (i *interpreter) visitTerm(e *expressionTerm) any {
	switch e.tokenType() {
	case PLUS:
		if i.evaluatesToString(e) {
			return fmt.Sprintf("%v%v", i.evaluate(e.expr()), i.evaluate(e.next()))
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

func (i *interpreter) visitLiteral(e *expressionLiteral) any {
	return e.value()
}

func (i *interpreter) visitGroup(e *expressionGroup) any {
	return i.evaluate(e.exprInterface)
}

func (i *interpreter) visitExpr(e *expression) any {
	return ""
}

func (i *interpreter) evaluatesToString(e exprInterface) bool {
	if e.expr() == nil || e.next() == nil {
		return e.tokenType() == STRING
	}
	return reflect.TypeOf(i.evaluate(e.expr())).Name() == "string" &&
		reflect.TypeOf(i.evaluate(e.next())).Name() == "string"
}

func (i *interpreter) parseFloat(e exprInterface) float64 {
	n := i.evaluate(e)
	switch n := n.(type) {
	case float64:
		return n
	}
	err := newError(fmt.Sprintf("Operand must be a number: %v", e.lexeme()), e.token().line)
	panic(err)
}

func (i *interpreter) hasSameType(a any, b any) bool {
	return reflect.TypeOf(a).Name() == reflect.TypeOf(b).Name()
}

func (i *interpreter) isTruthy(e exprInterface) bool {
	value := i.evaluate(e)
	switch value := value.(type) {
	case string:
		return value != ""
	case float64:
		return value != 0
	case bool:
		return value
	case nil:
		return false
	default:
		return false
	}
}

func (i *interpreter) stringify(val any) string {
	if val == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", val)
}
