package lox

type stmtVisitor interface {
	visitVarStmt(stmt *stmtVar)
	visitIfStmt(stmt *stmtIf)
	visitPrintStmt(stmt *stmtPrint)
	visitWhileStmt(stmt *stmtWhile)
	visitBlockStmt(stmt *stmtBlock)
	visitExprStmt(stmt *stmtExpr)
}
