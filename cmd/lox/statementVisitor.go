package lox

type stmtVisitor interface {
	visitPrintStmt(expr *stmtPrint)
	visitExprStmt(expr *stmtExpr)
}
