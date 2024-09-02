package lox

type callable interface {
	arity() int
	call(*interpreter, []any) any
}

type loxFunction struct {
	environment *environment
	declaration *stmtFun
}

type builtin struct {
	function func(*interpreter, []any) any
	lenArgs  int
}

func (f *loxFunction) String() string { return "<fn " + f.declaration.lexeme + ">" }
func (f *loxFunction) arity() int     { return len(f.declaration.params) }
func (f *loxFunction) call(i *interpreter, args []any) (value any) {
	defer func() {
		if r, ok := recover().(returnValue); ok {
			value = r.value
		}
	}()
	for i, param := range f.declaration.params {
		f.environment.define(param.lexeme, args[i])
	}
	block := f.declaration.stmtInterface.(*stmtBlock)
	i.executeBlock(block.statements, f.environment)
	return
}

func (b *builtin) String() string                      { return "<native fn>" }
func (b *builtin) arity() int                          { return b.lenArgs }
func (b *builtin) call(i *interpreter, args []any) any { return b.function(i, args) }
