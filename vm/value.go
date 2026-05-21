package main

type value float64

type valueArray []value

func (v *valueArray) init() {
	*v = make(valueArray, 0)
}

func (v *valueArray) write(b value) {
	*v = append(*v, b)
}
