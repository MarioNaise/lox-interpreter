package lox

type stmtVisitor interface {
	visitPrintStmt(expr *stmtPrint) string
	visitExprStmt(expr *stmtExpr) string
}
