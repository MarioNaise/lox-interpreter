package vm

import (
	"fmt"
	"os"
)

type interpretResult byte

const stackMax = 256

const (
	interpretOk interpretResult = iota
	interpretCompileError
	interpretRuntimeError
)

type vm struct {
	stack    [stackMax]value
	chunk    *chunk
	ip       int
	stackTop int
}

func (vm *vm) run() interpretResult {
	for {
		if debug {
			fmt.Printf("   stack: ")
			fmt.Println(vm.stack[:vm.stackTop])
			vm.chunk.disassembleInstruction(vm.ip)
		}
		instruction := vm.readByte()
		switch opCode(instruction) {
		case opConstant:
			constant := vm.readConstant()
			vm.push(constant)
		case opNil:
			vm.push(nilValue())
		case opTrue:
			vm.push(boolValue(true))
		case opFalse:
			vm.push(boolValue(false))
		case opEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(boolValue(a.equals(b)))
		case opAdd, opSubtract, opMultiply, opDivide, opGreater, opLess:
			if result := vm.binaryOp(opCode(instruction)); result != interpretOk {
				return result
			}
		case opNot:
			vm.push(boolValue(vm.pop().isFalsey()))
		case opNegate:
			if !vm.peek(0).isNumber() {
				vm.runtimeError("Operand must be a number.")
				return interpretRuntimeError
			}
			vm.push(numberValue(-vm.pop().value.(float64)))
		case opReturn:
			return interpretOk
		}
	}
}

func (vm *vm) interpret(src string) interpretResult {
	ch := currentChunk()
	if !compile(src, ch) {
		return interpretCompileError
	}

	vm.chunk = ch
	vm.ip = 0

	result := vm.run()

	return result
}

func (vm *vm) readByte() byte {
	b := vm.chunk.code[vm.ip]
	vm.ip++
	return b
}

func (vm *vm) readConstant() value {
	return vm.chunk.constants[vm.readByte()]
}

func (vm *vm) resetStack() {
	vm.stackTop = 0
}

func (vm *vm) runtimeError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	line := vm.chunk.lines[vm.ip-1]
	fmt.Fprintf(os.Stderr, "\n[line %d] in script\n", line)
	vm.resetStack()
}

func (vm *vm) push(value value) {
	if vm.stackTop >= stackMax {
		panic("Stack overflow")
	}
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *vm) pop() value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *vm) peek(distance int) value {
	return vm.stack[vm.stackTop-1-distance]
}

func (vm *vm) binaryOp(operator opCode) interpretResult {
	if !vm.peek(0).isNumber() || !vm.peek(1).isNumber() {
		vm.runtimeError("Operands must be numbers.")
		return interpretRuntimeError
	}
	b := vm.pop().value.(float64)
	a := vm.pop().value.(float64)
	switch operator {
	case opAdd:
		vm.push(numberValue(a + b))
	case opSubtract:
		vm.push(numberValue(a - b))
	case opMultiply:
		vm.push(numberValue(a * b))
	case opDivide:
		vm.push(numberValue(a / b))
	case opGreater:
		vm.push(boolValue(a > b))
	case opLess:
		vm.push(boolValue(a < b))
	}
	return interpretOk
}
