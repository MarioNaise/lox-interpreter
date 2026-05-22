package main

type opCode byte

const (
	opConstant opCode = iota
	opAdd
	opSubtract
	opMultiply
	opDivide
	opNegate
	opReturn
)

type chunk struct {
	code      []byte
	lines     []int
	constants valueArray
}

func (c *chunk) init() {
	c.code = make([]byte, 0)
	c.lines = make([]int, 0)
	c.constants.init()
}

func (c *chunk) write(b byte, line int) {
	c.code = append(c.code, b)
	c.lines = append(c.lines, line)
}

func (c *chunk) addConstant(p value) int {
	c.constants.write(p)
	return len(c.constants) - 1
}
