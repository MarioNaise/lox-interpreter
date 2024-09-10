package lox

type callable interface {
	arity() int
	call(*interpreter, []any, token) any
}

type loxClass struct {
	name string
}

type loxFunction struct {
	closure     *environment
	declaration *stmtFun
}

type builtin struct {
	function func(*interpreter, []any, token) any
	lenArgs  int
}

func (c *loxClass) String() string { return "<class " + c.name + ">" }

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
		f.closure.define(param.lexeme, args[i])
	}
	block := f.declaration.body.(*stmtBlock)
	i.executeBlock(block.statements, newEnvironment(f.closure))
	return
}

func (b *builtin) String() string                               { return "<native fn>" }
func (b *builtin) arity() int                                   { return b.lenArgs }
func (b *builtin) call(i *interpreter, args []any, t token) any { return b.function(i, args, t) }
