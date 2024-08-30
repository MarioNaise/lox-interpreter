package lox

type expressionVisitor interface {
	visitVar(expr *expressionVar) any
	visitAssignment(expr *expressionAssignment) any
	visitLogical(expr *expressionLogical) any
	visitEquality(expr *expressionEquality) any
	visitComparison(expr *expressionComparison) any
	visitTerm(expr *expressionTerm) any
	visitFactor(expr *expressionFactor) any
	visitUnary(expr *expressionUnary) any
	visitCall(expr *expressionCall) any
	visitLiteral(expr *expressionLiteral) any
	visitGroup(expr *expressionGroup) any
	visitExpr(expr *expression) any
}
