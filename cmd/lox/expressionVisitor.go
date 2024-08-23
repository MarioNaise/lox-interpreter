package lox

type expressionVisitor interface {
	visitEquality(expr *expressionEquality) any
	visitComparison(expr *expressionComparison) any
	visitTerm(expr *expressionTerm) any
	visitFactor(expr *expressionFactor) any
	visitUnary(expr *expressionUnary) any
	visitLiteral(expr *expressionLiteral) any
	visitGroup(expr *expressionGroup) any
	visitExpr(expr *expression) any
}
