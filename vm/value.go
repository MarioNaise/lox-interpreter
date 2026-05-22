package vm

type value float64

type valueArray []value

func (v *valueArray) write(b value) {
	*v = append(*v, b)
}
