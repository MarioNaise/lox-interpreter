package lox

type stmtInterface interface {
	accept(v stmtVisitor)
}

type stmtFun struct {
	body   stmtInterface
	name   token
	params []token
}

type stmtVar struct {
	initializer expression
	name        token
}

type stmtIf struct {
	condition  expression
	thenBranch stmtInterface
	elseBranch stmtInterface
}

type stmtReturn struct {
	value expression
}

type stmtWhile struct {
	condition expression
	body      stmtInterface
}

type stmtBlock struct {
	statements []stmtInterface
}

type stmtExpr struct {
	initializer expression
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

func (s *stmtReturn) accept(v stmtVisitor) {
	v.visitReturnStmt(s)
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
