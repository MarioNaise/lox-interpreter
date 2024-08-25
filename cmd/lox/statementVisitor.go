package lox

type stmtVisitor interface {
	visitVarStmt(stmt *stmtVar)
	visitPrintStmt(stmt *stmtPrint)
	visitExprStmt(stmt *stmtExpr)
}
