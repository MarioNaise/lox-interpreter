package lox

type stmt interface {
	accept(v stmtVisitor)
}

type stmtClass struct {
	name    token
	methods []stmtFun
}

type stmtFun struct {
	name   token
	body   stmt
	params []token
}

type stmtVar struct {
	initializer expression
	name        token
}

type stmtIf struct {
	condition  expression
	thenBranch stmt
	elseBranch stmt
}

type stmtReturn struct {
	value expression
	token
}

type stmtWhile struct {
	condition expression
	body      stmt
}

type stmtBlock struct {
	statements []stmt
}

type stmtExpr struct {
	initializer expression
}

func (s *stmtClass) accept(v stmtVisitor) {
	v.visitClassStmt(s)
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
