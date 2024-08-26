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

type expression struct {
	expression exprInterface
	right      exprInterface
	operator   token
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

func (e *expression) accept(v expressionVisitor) any {
	return v.visitExpr(e)
}

type expressionAssignment struct {
	exprInterface
}

func (e *expressionAssignment) accept(v expressionVisitor) any {
	return v.visitAssignment(e)
}

type expressionEquality struct {
	exprInterface
}

func (e *expressionEquality) accept(v expressionVisitor) any {
	return v.visitEquality(e)
}

type expressionComparison struct {
	exprInterface
}

func (e *expressionComparison) accept(v expressionVisitor) any {
	return v.visitComparison(e)
}

type expressionTerm struct {
	exprInterface
}

func (e *expressionTerm) accept(v expressionVisitor) any {
	return v.visitTerm(e)
}

type expressionFactor struct {
	exprInterface
}

func (e *expressionFactor) accept(v expressionVisitor) any {
	return v.visitFactor(e)
}

type expressionUnary struct {
	exprInterface
}

func (e *expressionUnary) accept(v expressionVisitor) any {
	return v.visitUnary(e)
}

type expressionLiteral struct {
	exprInterface
	val any
}

func (e *expressionLiteral) accept(v expressionVisitor) any {
	return v.visitLiteral(e)
}

func (e *expressionLiteral) value() any {
	return e.val
}

type expressionGroup struct {
	exprInterface
}

func (e *expressionGroup) accept(v expressionVisitor) any {
	return v.visitGroup(e)
}

type expressionVar struct {
	exprInterface
}

func (e *expressionVar) accept(v expressionVisitor) any {
	return v.visitVar(e)
}
