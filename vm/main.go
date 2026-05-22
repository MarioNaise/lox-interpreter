package main

func main() {
	c := new(chunk)
	constant := c.addConstant(2.2)
	c.write(byte(opConstant), 123)
	c.write(byte(constant), 123)
	constant = c.addConstant(3.4)
	c.write(byte(opConstant), 123)
	c.write(byte(constant), 123)

	c.write(byte(opAdd), 123)

	constant = c.addConstant(5.6)
	c.write(byte(opConstant), 123)
	c.write(byte(constant), 123)

	c.write(byte(opDivide), 123)
	c.write(byte(opNegate), 123)
	c.write(byte(opReturn), 123)
	vm := new(vm)
	vm.init()
	vm.interpret(c)
}
