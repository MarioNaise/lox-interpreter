package lox

type stmtVisitor interface {
	visitVarStmt(stmt *stmtVar)
	visitIfStmt(stmt *stmtIf)
	visitPrintStmt(stmt *stmtPrint)
	visitBlockStmt(stmt *stmtBlock)
	visitExprStmt(stmt *stmtExpr)
}
