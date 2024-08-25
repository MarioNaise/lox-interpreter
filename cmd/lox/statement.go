package lox

type stmt struct {
	exprInterface
}

type stmtInterface interface {
	expr() exprInterface
	getName() token
	accept(v stmtVisitor)
}

type stmtExpr struct {
	initializer exprInterface
	token
}

func (s *stmtExpr) expr() exprInterface {
	return s.initializer
}

func (s *stmtExpr) getName() token {
	return s.token
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

type stmtVar struct {
	stmtInterface
}

func (s *stmtVar) accept(v stmtVisitor) {
	v.visitVarStmt(s)
}
