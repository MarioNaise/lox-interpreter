package lox

type stmtVisitor interface {
	visitVarStmt(stmt *stmtVar)
	visitPrintStmt(stmt *stmtPrint)
	visitBlockStmt(stmt *stmtBlock)
	visitExprStmt(stmt *stmtExpr)
}
