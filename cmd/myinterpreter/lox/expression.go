package lox

type exprInterface interface {
	get() *expression
	accept(v Visitor) string
}

type expression struct {
	expression exprInterface
	right      exprInterface
	operator   token
}

func (e *expression) get() *expression {
	return e
}

func (e *expression) accept(v Visitor) string {
	return v.visitExpr(e)
}

type expressionLiteral struct {
	exprInterface
}

func (e *expressionLiteral) accept(v Visitor) string {
	return v.visitLiteral(e)
}

type expressionGroup struct {
	exprInterface
}

func (e *expressionGroup) accept(v Visitor) string {
	return v.visitGroup(e)
}

type expressionEquality struct {
	exprInterface
}

func (e *expressionEquality) accept(v Visitor) string {
	return v.visitEquality(e)
}

type expressionComparison struct {
	exprInterface
}

func (e *expressionComparison) accept(v Visitor) string {
	return v.visitComparison(e)
}

type expressionTerm struct {
	exprInterface
}

func (e *expressionTerm) accept(v Visitor) string {
	return v.visitTerm(e)
}

type expressionFactor struct {
	exprInterface
}

func (e *expressionFactor) accept(v Visitor) string {
	return v.visitFactor(e)
}

type expressionUnary struct {
	exprInterface
}

func (e *expressionUnary) accept(v Visitor) string {
	return v.visitUnary(e)
}
