package lox

type callableInterface interface {
	arity() int
	call(*interpreter, []any) any
}

type loxFunction struct {
	environment *environment
	declaration *stmtFun
}

func (f *loxFunction) arity() int { return len(f.declaration.params) }
func (f *loxFunction) call(i *interpreter, args []any) any {
	f.environment = newEnvironment(i.environment)
	for i, param := range f.declaration.params {
		f.environment.define(param.lexeme, args[i])
	}
	block := f.declaration.stmtInterface.(*stmtBlock)
	i.executeBlock(block.statements, f.environment)
	return nil
}

type callable struct {
	String   func() string
	evaluate func() any
	args     []any
}

func (c *callable) arity() int                       { return len(c.args) }
func (c *callable) call(_ *interpreter, _ []any) any { return c.evaluate() }
