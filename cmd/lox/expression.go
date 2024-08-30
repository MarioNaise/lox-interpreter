package lox

type exprInterface interface {
	accept(v expressionVisitor) any
	expr() exprInterface
	next() exprInterface
	token() token
	tokenType() string
	lexeme() string
	literal() string
}

type expressionVar struct {
	exprInterface
}

type expressionAssignment struct {
	exprInterface
}

type expressionLogical struct {
	exprInterface
}

type expressionEquality struct {
	exprInterface
}

type expressionComparison struct {
	exprInterface
}

type expressionTerm struct {
	exprInterface
}

type expressionFactor struct {
	exprInterface
}

type expressionUnary struct {
	exprInterface
}

type expressionCall struct {
	exprInterface
	callee exprInterface
	args   []exprInterface
}

type expressionLiteral struct {
	exprInterface
	val any
}

type expressionGroup struct {
	exprInterface
}

type expression struct {
	expression exprInterface
	right      exprInterface
	operator   token
}

func (e *expressionVar) accept(v expressionVisitor) any {
	return v.visitVar(e)
}

func (e *expressionAssignment) accept(v expressionVisitor) any {
	return v.visitAssignment(e)
}

func (e *expressionLogical) accept(v expressionVisitor) any {
	return v.visitLogical(e)
}

func (e *expressionEquality) accept(v expressionVisitor) any {
	return v.visitEquality(e)
}

func (e *expressionComparison) accept(v expressionVisitor) any {
	return v.visitComparison(e)
}

func (e *expressionTerm) accept(v expressionVisitor) any {
	return v.visitTerm(e)
}

func (e *expressionFactor) accept(v expressionVisitor) any {
	return v.visitFactor(e)
}

func (e *expressionCall) accept(v expressionVisitor) any {
	return v.visitCall(e)
}

func (e *expressionUnary) accept(v expressionVisitor) any {
	return v.visitUnary(e)
}

func (e *expressionLiteral) accept(v expressionVisitor) any {
	return v.visitLiteral(e)
}

func (e *expressionGroup) accept(v expressionVisitor) any {
	return v.visitGroup(e)
}

func (e *expression) accept(v expressionVisitor) any {
	return v.visitExpr(e)
}

func (e *expressionLiteral) value() any {
	return e.val
}

func (e *expression) expr() exprInterface {
	return e.expression
}

func (e *expression) next() exprInterface {
	return e.right
}

func (e *expression) token() token {
	return e.operator
}

func (e *expression) tokenType() string {
	return e.operator.tokenType
}

func (e *expression) lexeme() string {
	return e.operator.lexeme
}

func (e *expression) literal() string {
	return e.operator.literal
}
