package lox

type stmt struct {
	exprInterface
}

type stmtInterface interface {
	expr() exprInterface
	name() token
	accept(v stmtVisitor)
}

type stmtExpr struct {
	initializer exprInterface
	token
}

func (s *stmtExpr) expr() exprInterface {
	return s.initializer
}

func (s *stmtExpr) name() token {
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

type stmtIf struct {
	stmtInterface
	condition  exprInterface
	thenBranch stmtInterface
	elseBranch stmtInterface
}

func (s *stmtIf) accept(v stmtVisitor) {
	v.visitIfStmt(s)
}

type stmtVar struct {
	stmtInterface
}

func (s *stmtVar) accept(v stmtVisitor) {
	v.visitVarStmt(s)
}

type stmtBlock struct {
	stmtInterface
	statements []stmtInterface
}

func (s *stmtBlock) accept(v stmtVisitor) {
	v.visitBlockStmt(s)
}
