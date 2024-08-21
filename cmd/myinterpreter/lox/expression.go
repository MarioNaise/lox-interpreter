package lox

type exprInterface interface {
	accept(v expressionVisitor) string
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

func (e *expression) token() token {
	return e.operator
}

func (e *expression) expr() exprInterface {
	return e.expression
}

func (e *expression) next() exprInterface {
	return e.right
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

func (e *expression) accept(v expressionVisitor) string {
	return v.visitExpr(e)
}

type expressionLiteral struct {
	exprInterface
}

func (e *expressionLiteral) accept(v expressionVisitor) string {
	return v.visitLiteral(e)
}

type expressionGroup struct {
	exprInterface
}

func (e *expressionGroup) accept(v expressionVisitor) string {
	return v.visitGroup(e)
}

type expressionEquality struct {
	exprInterface
}

func (e *expressionEquality) accept(v expressionVisitor) string {
	return v.visitEquality(e)
}

type expressionComparison struct {
	exprInterface
}

func (e *expressionComparison) accept(v expressionVisitor) string {
	return v.visitComparison(e)
}

type expressionTerm struct {
	exprInterface
}

func (e *expressionTerm) accept(v expressionVisitor) string {
	return v.visitTerm(e)
}

type expressionFactor struct {
	exprInterface
}

func (e *expressionFactor) accept(v expressionVisitor) string {
	return v.visitFactor(e)
}

type expressionUnary struct {
	exprInterface
}

func (e *expressionUnary) accept(v expressionVisitor) string {
	return v.visitUnary(e)
}
