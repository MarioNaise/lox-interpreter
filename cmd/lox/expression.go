package lox

type expression interface {
	accept(v expressionVisitor) any
	expr() expression
	next() expression
	token() token
	tokenType() string
	lexeme() string
	literal() string
}

type expressionVar struct {
	expression
}

type expressionAssignment struct {
	expression
}

type expressionLogical struct {
	expression
}

type expressionEquality struct {
	expression
}

type expressionComparison struct {
	expression
}

type expressionTerm struct {
	expression
}

type expressionFactor struct {
	expression
}

type expressionUnary struct {
	expression
}

type expressionGet struct {
	expression
	name token
}

type expressionCall struct {
	expression
	args []expression
}

type expressionLiteral struct {
	expression
	val any
}

type expressionGroup struct {
	expression
}

type exp struct {
	expression expression
	right      expression
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

func (e *expressionGet) accept(v expressionVisitor) any {
	return v.visitGet(e)
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

func (e *exp) accept(v expressionVisitor) any {
	return v.visitExpr(e)
}

func (e *expressionLiteral) value() any {
	return e.val
}

func (e *exp) expr() expression {
	return e.expression
}

func (e *exp) next() expression {
	return e.right
}

func (e *exp) token() token {
	return e.operator
}

func (e *exp) tokenType() string {
	return e.operator.tokenType
}

func (e *exp) lexeme() string {
	return e.operator.lexeme
}

func (e *exp) literal() string {
	return e.operator.literal
}
