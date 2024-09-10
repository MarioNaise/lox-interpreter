package lox

import (
	"container/list"
	"fmt"
)

type fnType int

const (
	none fnType = iota
	function
)

type resolver struct {
	interpreter *interpreter
	scopes      *list.List
	currentFun  fnType
}

func newResolver(i *interpreter) *resolver {
	r := resolver{i, list.New(), none}
	return &r
}

func (r *resolver) resolve(stmts []stmt) {
	for _, s := range stmts {
		r.resolveStmt(s)
	}
}

func (r *resolver) resolveStmt(s stmt) {
	if s != nil {
		s.accept(r)
	}
}

func (r *resolver) resolveExpr(e expression) {
	if e != nil {
		e.accept(r)
	}
}

func (r *resolver) visitClassStmt(stmt *stmtClass) {
	r.declare(stmt.name)
	r.define(stmt.name)
}

func (r *resolver) visitFunStmt(stmt *stmtFun) {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt)
}

func (r *resolver) resolveFunction(stmt *stmtFun) {
	enclosingFn := r.currentFun
	defer func() { r.currentFun = enclosingFn }()
	r.currentFun = function
	r.beginScope()
	for _, param := range stmt.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmt(stmt.body)
	r.endScope()
}

func (r *resolver) visitVarStmt(stmt *stmtVar) {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpr(stmt.initializer)
	}
	r.define(stmt.name)
}

func (r *resolver) declare(name token) {
	if r.scopes.Len() == 0 {
		return
	}
	scope := r.scopes.Back().Value.(map[string]bool)
	if _, ok := scope[name.lexeme]; ok {
		err := newError(fmt.Sprintf("Identifier '%s' already declared in this scope.", name.lexeme), name.line)
		panic(err)
	}
	scope[name.lexeme] = false
}

func (r *resolver) define(name token) {
	if r.scopes.Len() == 0 {
		return
	}
	scope := r.scopes.Back().Value.(map[string]bool)
	scope[name.lexeme] = true
}

func (r *resolver) visitIfStmt(stmt *stmtIf) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
}

func (r *resolver) visitReturnStmt(stmt *stmtReturn) {
	if r.currentFun == none {
		err := newError("Can't return from top-level code.", stmt.line)
		panic(err)
	}
	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
}

func (r *resolver) visitWhileStmt(stmt *stmtWhile) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
}

func (r *resolver) visitBlockStmt(stmt *stmtBlock) {
	r.beginScope()
	r.resolve(stmt.statements)
	r.endScope()
}

func (r *resolver) beginScope() {
	r.scopes.PushBack(map[string]bool{})
}

func (r *resolver) endScope() {
	r.scopes.Remove(r.scopes.Back())
}

func (r *resolver) visitExprStmt(stmt *stmtExpr) {
	r.resolveExpr(stmt.initializer)
}

func (r *resolver) visitVar(expr *expressionVar) any {
	if r.scopes.Len() > 0 {
		hasValue, ok := r.scopes.Back().Value.(map[string]bool)[expr.lexeme()]
		if r.scopes.Len() > 0 && ok && !hasValue {
			err := newError(fmt.Sprintf("Can't access '%s' in its own initializer.", expr.lexeme()), expr.token().line)
			panic(err)
		}
	}
	r.resolveLocal(expr)
	return nil
}

func (r *resolver) resolveLocal(expr expression) {
	for i := r.scopes.Len() - 1; i >= 0; i-- {
		scope := getNthOfList(r.scopes, i).Value.(map[string]bool)
		if _, ok := scope[expr.lexeme()]; ok {
			r.interpreter.resolve(expr, r.scopes.Len()-1-i)
			return
		}
	}
}

func getNthOfList(l *list.List, n int) *list.Element {
	e := l.Front()
	for i := 0; i < n; i++ {
		e = e.Next()
	}
	return e
}

func (r *resolver) visitAssignment(expr *expressionAssignment) any {
	r.resolveExpr(expr.next())
	r.resolveLocal(expr)
	return nil
}

func (r *resolver) visitLogical(expr *expressionLogical) any {
	return r.defaultResolver(expr)
}

func (r *resolver) visitEquality(expr *expressionEquality) any {
	return r.defaultResolver(expr)
}

func (r *resolver) visitComparison(expr *expressionComparison) any {
	return r.defaultResolver(expr)
}

func (r *resolver) visitTerm(expr *expressionTerm) any {
	return r.defaultResolver(expr)
}

func (r *resolver) visitFactor(expr *expressionFactor) any {
	return r.defaultResolver(expr)
}

func (r *resolver) visitUnary(expr *expressionUnary) any {
	r.resolveExpr(expr.next())
	return nil
}

func (r *resolver) visitCall(expr *expressionCall) any {
	r.resolveExpr(expr.expression)
	for _, arg := range expr.args {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *resolver) visitLiteral(expr *expressionLiteral) any { return nil }

func (r *resolver) visitGroup(expr *expressionGroup) any {
	r.resolveExpr(expr.expression)
	return nil
}

func (r *resolver) visitExpr(expr *exp) any { return nil }

func (r *resolver) defaultResolver(expr expression) any {
	r.resolveExpr(expr.expr())
	r.resolveExpr(expr.next())
	return nil
}
