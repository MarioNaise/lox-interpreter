package lox

type stmtVisitor interface {
	visitVarStmt(stmt *stmtVar)
	visitPrintStmt(stmt *stmtPrint)
	visitIfStmt(stmt *stmtIf)
	visitBlockStmt(stmt *stmtBlock)
	visitExprStmt(stmt *stmtExpr)
}
