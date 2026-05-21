package main

func main() {
	ch := new(chunk)
	const1 := ch.addConstant(1.2)
	ch.write(byte(opConstant), 123)
	ch.write(byte(const1), 123)
	ch.write(byte(opReturn), 123)
	ch.disassemble("test chunk")
}
