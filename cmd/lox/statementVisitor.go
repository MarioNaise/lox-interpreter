package lox

type stmtVisitor interface {
	visitFunStmt(stmt *stmtFun)
	visitVarStmt(stmt *stmtVar)
	visitIfStmt(stmt *stmtIf)
	visitReturnStmt(stmt *stmtReturn)
	visitWhileStmt(stmt *stmtWhile)
	visitBlockStmt(stmt *stmtBlock)
	visitExprStmt(stmt *stmtExpr)
}
