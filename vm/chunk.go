package vm

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

func newChunk() *chunk {
	c := new(chunk)
	c.code = make([]byte, 0)
	c.lines = make([]int, 0)
	c.constants = make(valueArray, 0)
	return c
}

func (c *chunk) write(b byte, line int) {
	c.code = append(c.code, b)
	c.lines = append(c.lines, line)
}

func (c *chunk) addConstant(p value) int {
	c.constants.write(p)
	return len(c.constants) - 1
}
