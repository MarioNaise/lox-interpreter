package lox

type stmt struct {
	exprInterface
}

type stmtInterface interface {
	accept(v stmtVisitor) string
}

type stmtExpr struct {
	expression exprInterface
}

func (e *stmtExpr) accept(v stmtVisitor) string {
	return v.visitExprStmt(e)
}

type stmtPrint struct {
	expression exprInterface
	stmtInterface
}

func (e *stmtPrint) accept(v stmtVisitor) string {
	return v.visitPrintStmt(e)
}
