package lox

type callable struct {
	String   func() string
	evaluate func() any
	args     []any
}

func (c *callable) arity() int                       { return len(c.args) }
func (c *callable) call(_ *interpreter, _ []any) any { return c.evaluate() }
