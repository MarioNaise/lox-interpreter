package lox

import (
	"fmt"
	"strings"
)

const (
	BLOCK     = "BLOCK"
	BLOCK_END = "BLOCK_END"
	GROUP     = "GROUP"
)

type astPrinter struct{}

func (a *astPrinter) print(stmts []stmt) {
	for _, s := range stmts {
		s.accept(a)
	}
}

func (a *astPrinter) printExpr(e expression) {
	if e == nil {
		return
	}
	fmt.Println(e.accept(a))
}

func (a *astPrinter) visitClassStmt(stmt *stmtClass) {
	fmt.Println(CLASS + ":" + stmt.name.lexeme)
	for _, method := range stmt.methods {
		method.accept(a)
	}
	fmt.Println(CLASS + "_END")
}

func (a *astPrinter) visitFunStmt(s *stmtFun) {
	a.prefix(fmt.Sprintf("%s:%s", FUN, s.name.lexeme))
	p := []string{}
	for _, param := range s.params {
		p = append(p, param.lexeme)
	}
	fmt.Println("[", strings.Join(p, ", "), "]")
	s.body.accept(a)
}

func (a *astPrinter) visitVarStmt(s *stmtVar) {
	a.prefix(VAR + ":" + s.name.lexeme)
	a.printExpr(s.initializer)
}

func (a *astPrinter) visitIfStmt(s *stmtIf) {
	a.prefix(IF)
	a.printExpr(s.condition)
	s.thenBranch.accept(a)
	if s.elseBranch != nil {
		fmt.Println(ELSE)
		s.elseBranch.accept(a)
	}
}

func (a *astPrinter) visitReturnStmt(s *stmtReturn) {
	a.prefix(RETURN)
	a.printExpr(s.value)
}

func (a *astPrinter) visitWhileStmt(s *stmtWhile) {
	a.prefix(WHILE)
	a.printExpr(s.condition)
	s.body.accept(a)
}

func (a *astPrinter) visitBlockStmt(s *stmtBlock) {
	fmt.Println(BLOCK)
	for _, stmt := range s.statements {
		stmt.accept(a)
	}
	fmt.Println(BLOCK_END)
}

func (a *astPrinter) visitExprStmt(s *stmtExpr) {
	a.printExpr(s.initializer)
}

func (a *astPrinter) visitVar(e *expressionVar) any {
	return fmt.Sprintf("%s %s", VAR, e.lexeme())
}

func (a *astPrinter) visitAssignment(e *expressionAssignment) any {
	return fmt.Sprintf("%s:%s %v", VAR, e.expr().lexeme(), e.next().accept(a))
}

func (a *astPrinter) visitSet(e *expressionSet) any {
	return fmt.Sprintf("%v.%v: %v", e.expression.accept(a), e.name.lexeme, e.value.accept(a))
}

func (a *astPrinter) visitLogical(e *expressionLogical) any {
	return a.parenthesized(strings.ToUpper(e.lexeme()), e.expr(), e.next())
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

func (a *astPrinter) visitGet(e *expressionGet) any {
	return fmt.Sprintf("%v.%v", e.expression.accept(a), e.name.lexeme)
}

func (a *astPrinter) visitCall(e *expressionCall) any {
	argStr := a.joinExprs(e.args)
	return fmt.Sprintf("(%s [%s])", e.lexeme(), argStr)
}

func (a *astPrinter) visitLiteral(e *expressionLiteral) any {
	return a.primary(e)
}

func (a *astPrinter) visitGroup(e *expressionGroup) any {
	return a.parenthesized(GROUP, e.expression)
}

func (a *astPrinter) visitExpr(e *exp) any { return "" }

func (a *astPrinter) primary(e expression) any {
	if e.literal() != NULL {
		return e.literal()
	}
	return e.lexeme()
}

func (a *astPrinter) parenthesized(name string, e ...expression) string {
	joined := a.joinExprs(e)
	return fmt.Sprintf("(%s %s)", name, joined)
}

func (a *astPrinter) joinExprs(e []expression) string {
	s := []string{}
	for _, expr := range e {
		if expr != nil {
			s = append(s, expr.accept(a).(string))
		}
	}
	joined := strings.Join(s, " ")
	return joined
}

func (a *astPrinter) defaultString(e expression) string {
	return a.parenthesized(e.lexeme(), e.expr(), e.next())
}

func (a *astPrinter) prefix(s string) {
	fmt.Print(s, " ")
}
