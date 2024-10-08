package lox

import (
	"fmt"
	"reflect"
)

type interpreter struct {
	*resolver
	*parser
	*environment
	locals map[expression]int
	index  string
}

func newInterpreter(str string, index string) *interpreter {
	glob := newEnvironment(nil)
	glob.values = globals()
	locals := make(map[expression]int)
	p := newParser(str)
	i := interpreter{nil, p, glob, locals, index}
	i.resolver = newResolver(&i)
	return &i
}

func (i *interpreter) evaluate(e expression) any {
	return e.accept(i)
}

func (i *interpreter) execute(s stmt) {
	defer i.syncOnError()
	s.accept(i)
}

func (i *interpreter) interpret(stmts []stmt) {
	for _, s := range stmts {
		i.execute(s)
	}
}

func (i *interpreter) resolve(expr expression, depth int) {
	i.locals[expr] = depth
}

func (i *interpreter) visitClassStmt(stmt *stmtClass) {
	i.environment.define(stmt.name.lexeme, nil)
	methods := make(map[string]*loxFunction)
	for _, m := range stmt.methods {
		fun := &loxFunction{newEnvironment(i.environment), m}
		methods[m.name.lexeme] = fun
	}
	class := &loxClass{methods, stmt.name.lexeme}
	i.assign(stmt.name, class)
}

func (i *interpreter) visitFunStmt(s *stmtFun) {
	function := &loxFunction{newEnvironment(i.environment), s}
	i.environment.define(s.name.lexeme, function)
}

func (i *interpreter) visitVarStmt(s *stmtVar) {
	var val any
	name := s.name.lexeme
	if s.initializer != nil {
		val = i.evaluate(s.initializer)
	}
	i.environment.define(name, val)
}

func (i *interpreter) visitIfStmt(s *stmtIf) {
	if i.isTruthy(s.condition) {
		i.execute(s.thenBranch)
	} else if s.elseBranch != nil {
		i.execute(s.elseBranch)
	}
}

type returnValue struct {
	value any
}

func (i *interpreter) visitReturnStmt(s *stmtReturn) {
	if s.value == nil {
		// TODO: handle return globally and in REPL expressions
		panic(returnValue{nil})
	}
	panic(returnValue{i.evaluate(s.value)})
}

func (i *interpreter) visitWhileStmt(s *stmtWhile) {
	for i.isTruthy(s.condition) {
		i.execute(s.body)
	}
}

func (i *interpreter) visitBlockStmt(s *stmtBlock) {
	i.executeBlock(s.statements, newEnvironment(i.environment))
}

func (i *interpreter) executeBlock(stmts []stmt, env *environment) {
	prevEnv := i.environment
	defer func() { i.environment = prevEnv }()
	i.environment = env
	i.interpret(stmts)
}

func (i *interpreter) visitExprStmt(s *stmtExpr) {
	i.evaluate(s.initializer)
}

func (i *interpreter) visitVar(e *expressionVar) any {
	return i.lookupVariable(e)
}

func (i *interpreter) lookupVariable(e expression) any {
	distance, ok := i.locals[e]
	if ok {
		return i.getAt(distance, e.token())
	}
	return i.get(e.token())
}

func (i *interpreter) visitAssignment(e *expressionAssignment) any {
	value := i.evaluate(e.next())
	distance, ok := i.locals[e]
	if ok {
		i.assignAt(distance, e.expr().token(), value)
	} else {
		i.assign(e.expr().token(), value)
	}
	return value
}

func (i *interpreter) visitSet(expr *expressionSet) any {
	object := i.evaluate(expr.expression)
	instance, ok := object.(*loxInstance)
	if !ok {
		err := newError("Only instances have fields.", expr.token().line)
		panic(err)
	}
	val := i.evaluate(expr.value)
	instance.set(expr.name, val)
	return val
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
		if ok, left, right := i.evaluatesToString(e); ok {
			return fmt.Sprintf("%v%v", left, right)
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

func (i *interpreter) visitGet(expr *expressionGet) any {
	object := i.evaluate(expr.expression)
	instance, ok := object.(*loxInstance)
	if !ok {
		err := newError("Only instances have properties.", expr.token().line)
		panic(err)
	}
	return instance.get(expr.name)
}

func (i *interpreter) visitCall(e *expressionCall) any {
	defer recoverLoxError(e.token())
	callee := i.evaluate(e.expression)
	args := make([]any, 0)
	for _, arg := range e.args {
		args = append(args, i.evaluate(arg))
	}
	function, ok := callee.(callable)
	if !ok {
		panic(newError("Can only call functions and classes.", e.token().line))
	}
	if len(e.args) != function.arity() {
		err := newError(fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(e.args)), e.token().line)
		panic(err)
	}
	return function.call(i, args, e.token())
}

func (i *interpreter) visitLiteral(e *expressionLiteral) any {
	return e.value()
}

func (i *interpreter) visitGroup(e *expressionGroup) any {
	return i.evaluate(e.expression)
}

func (i *interpreter) visitExpr(e *exp) any {
	return ""
}

func (i *interpreter) evaluatesToString(e expression) (bool, any, any) {
	left := i.evaluate(e.expr())
	right := i.evaluate(e.next())
	if left == nil {
		left = "nil"
	}
	if right == nil {
		right = "nil"
	}
	return reflect.TypeOf(left).Name() == "string" &&
		reflect.TypeOf(right).Name() == "string", left, right
}

func (i *interpreter) parseFloat(e expression) float64 {
	n := i.evaluate(e)
	switch n := n.(type) {
	case float64:
		return n
	}
	err := newError(fmt.Sprintf("Operand must be a number: %v", e.lexeme()), e.token().line)
	panic(err)
}

func (i *interpreter) hasSameType(a any, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	return reflect.TypeOf(a).Name() == reflect.TypeOf(b).Name()
}

func (i *interpreter) isTruthy(e expression) bool {
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

func (i *interpreter) syncOnError() {
	if r := recover(); r != nil {
		i.synchronize()
		panic(r)
	}
}
