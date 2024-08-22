package lox

type stmt struct {
	exprInterface
}

type stmtInterface interface {
	expr() exprInterface
	accept(v stmtVisitor)
}

type stmtExpr struct {
	expression exprInterface
}

func (s *stmtExpr) expr() exprInterface {
	return s.expression
}

func (s *stmtExpr) accept(v stmtVisitor) {
	v.visitExprStmt(s)
}

type stmtPrint struct {
	stmtInterface
}

func (s *stmtPrint) accept(v stmtVisitor) {
	v.visitPrintStmt(s)
}
