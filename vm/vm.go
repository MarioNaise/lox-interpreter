package main

import (
	"fmt"
)

type interpretResult byte

const stackMax = 256

const (
	interpretOk interpretResult = iota
	interpretCompileError
	interpretRuntimeError
)

type vm struct {
	ip       int
	stackTop int
	chunk    *chunk
	stack    [stackMax]value
}

func (vm *vm) init() {
	vm.chunk = new(chunk)
	vm.chunk.init()
}

func (vm *vm) interpret(chunk *chunk) interpretResult {
	vm.chunk = chunk
	return vm.run()
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

func (vm *vm) binaryOp(operator opCode) {
	b := vm.pop()
	a := vm.pop()
	var result value
	switch operator {
	case opAdd:
		result = a + b
	case opSubtract:
		result = a - b
	case opMultiply:
		result = a * b
	case opDivide:
		result = a / b
	}
	vm.push(result)
}

func (vm *vm) run() interpretResult {
	for {
		if debug {
			fmt.Printf("          ")
			fmt.Printf("%v\n", vm.stack[:vm.stackTop])
			vm.chunk.disassembleInstruction(vm.ip)
		}
		instruction := vm.readByte()
		switch opCode(instruction) {
		case opConstant:
			constant := vm.readConstant()
			fmt.Println("opConstant", constant)
			vm.push(constant)
		case opAdd, opSubtract, opMultiply, opDivide:
			vm.binaryOp(opCode(instruction))
		case opNegate:
			vm.push(-vm.pop())
		case opReturn:
			fmt.Println("opReturn", vm.pop())
			return interpretOk
		}
	}
}
