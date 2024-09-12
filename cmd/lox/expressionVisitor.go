package lox

type expressionVisitor interface {
	visitVar(expr *expressionVar) any
	visitAssignment(expr *expressionAssignment) any
	visitSet(expr *expressionSet) any
	visitLogical(expr *expressionLogical) any
	visitEquality(expr *expressionEquality) any
	visitComparison(expr *expressionComparison) any
	visitTerm(expr *expressionTerm) any
	visitFactor(expr *expressionFactor) any
	visitUnary(expr *expressionUnary) any
	visitGet(expr *expressionGet) any
	visitCall(expr *expressionCall) any
	visitLiteral(expr *expressionLiteral) any
	visitGroup(expr *expressionGroup) any
	visitExpr(expr *exp) any
}
