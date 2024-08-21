package lox

type expressionVisitor interface {
	visitEquality(expr *expressionEquality) string
	visitComparison(expr *expressionComparison) string
	visitTerm(expr *expressionTerm) string
	visitFactor(expr *expressionFactor) string
	visitUnary(expr *expressionUnary) string
	visitLiteral(expr *expressionLiteral) string
	visitGroup(expr *expressionGroup) string
	visitExpr(expr *expression) string
}
