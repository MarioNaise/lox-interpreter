package main

import "fmt"

func (c *chunk) disassemble(s string) {
	fmt.Printf("== %s ==\n", s)

	for offset := 0; offset < len(c.code); {
		offset = c.disassembleInstruction(offset)
	}
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
