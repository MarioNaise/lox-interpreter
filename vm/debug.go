package vm

import "fmt"

const debug = true

func (c *chunk) disassemble(s string) {
	fmt.Printf("== %s ==\n", s)

	for offset := 0; offset < len(c.code); {
		offset = c.disassembleInstruction(offset)
	}
	fmt.Println("========")
}

func (c *chunk) disassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 &&
		c.lines[offset] == c.lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", c.lines[offset])
	}

	instruction := (c.code)[offset]
	switch opCode(instruction) {
	case opConstant:
		return c.constantInstruction("OP_CONSTANT", offset)
	case opNil:
		return c.simpleInstruction("OP_NIL", offset)
	case opTrue:
		return c.simpleInstruction("OP_TRUE", offset)
	case opFalse:
		return c.simpleInstruction("OP_FALSE", offset)
	case opEqual:
		return c.simpleInstruction("OP_EQUAL", offset)
	case opGreater:
		return c.simpleInstruction("OP_GREATER", offset)
	case opLess:
		return c.simpleInstruction("OP_LESS", offset)
	case opAdd:
		return c.simpleInstruction("OP_ADD", offset)
	case opSubtract:
		return c.simpleInstruction("OP_SUBTRACT", offset)
	case opMultiply:
		return c.simpleInstruction("OP_MULTIPLY", offset)
	case opDivide:
		return c.simpleInstruction("OP_DIVIDE", offset)
	case opNot:
		return c.simpleInstruction("OP_NOT", offset)
	case opNegate:
		return c.simpleInstruction("OP_NEGATE", offset)
	case opReturn:
		return c.simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func (c *chunk) constantInstruction(name string, offset int) int {
	constant := (c.code)[offset+1]
	fmt.Printf("%-16s %4d '%v'\n", name, constant, c.constants[constant])
	return offset + 2
}

func (c *chunk) simpleInstruction(name string, offset int) int {
	fmt.Println(name)
	return offset + 1
}
