package lox

type callable interface {
	arity() int
	call(*interpreter, []any, token) any
}

type loxFunction struct {
	environment *environment
	declaration *stmtFun
}

type builtin struct {
	function func(*interpreter, []any, token) any
	lenArgs  int
}

func (f *loxFunction) String() string { return "<fn " + f.declaration.name.lexeme + ">" }
func (f *loxFunction) arity() int     { return len(f.declaration.params) }
func (f *loxFunction) call(i *interpreter, args []any, t token) (value any) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case returnValue:
				value = r.value
				return
			default:
				panic(r)
			}
		}
	}()
	for i, param := range f.declaration.params {
		f.environment.define(param.lexeme, args[i])
	}
	block := f.declaration.body.(*stmtBlock)
	i.executeBlock(block.statements, f.environment)
	return
}

func (b *builtin) String() string                               { return "<native fn>" }
func (b *builtin) arity() int                                   { return b.lenArgs }
func (b *builtin) call(i *interpreter, args []any, t token) any { return b.function(i, args, t) }
