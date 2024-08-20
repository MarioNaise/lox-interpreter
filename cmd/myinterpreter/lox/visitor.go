package lox

type Visitor interface {
	visitExpr(expr *expression) string
	visitEquality(expr *expressionEquality) string
	visitComparison(expr *expressionComparison) string
	visitTerm(expr *expressionTerm) string
	visitFactor(expr *expressionFactor) string
	visitUnary(expr *expressionUnary) string
	visitLiteral(expr *expressionLiteral) string
	visitGroup(expr *expressionGroup) string
}
