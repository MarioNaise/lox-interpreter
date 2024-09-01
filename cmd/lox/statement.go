package lox

type stmtInterface interface {
	expr() exprInterface
	name() token
	accept(v stmtVisitor)
}

type stmtFun struct {
	stmtInterface
	token
	params []token
}

type stmtVar struct {
	stmtInterface
}

type stmtIf struct {
	stmtInterface
	condition  exprInterface
	thenBranch stmtInterface
	elseBranch stmtInterface
}

type stmtPrint struct {
	stmtInterface
}

type stmtWhile struct {
	stmtInterface
	condition exprInterface
	body      stmtInterface
}

type stmtBlock struct {
	stmtInterface
	statements []stmtInterface
}

type stmtExpr struct {
	initializer exprInterface
	token
}

func (s *stmtFun) accept(v stmtVisitor) {
	v.visitFunStmt(s)
}

func (s *stmtVar) accept(v stmtVisitor) {
	v.visitVarStmt(s)
}

func (s *stmtIf) accept(v stmtVisitor) {
	v.visitIfStmt(s)
}

func (s *stmtPrint) accept(v stmtVisitor) {
	v.visitPrintStmt(s)
}

func (s *stmtWhile) accept(v stmtVisitor) {
	v.visitWhileStmt(s)
}

func (s *stmtBlock) accept(v stmtVisitor) {
	v.visitBlockStmt(s)
}

func (s *stmtExpr) accept(v stmtVisitor) {
	v.visitExprStmt(s)
}

func (s *stmtExpr) expr() exprInterface {
	return s.initializer
}

func (s *stmtExpr) name() token {
	return s.token
}
