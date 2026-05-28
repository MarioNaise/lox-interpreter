package vm

import "fmt"

type valueType int

const (
	valBool valueType = iota
	valNil
	valNumber
	valString
)

type value struct {
	valueType
	value any
}

func (v value) isBool() bool {
	return v.valueType == valBool
}

func (v value) isNil() bool {
	return v.valueType == valNil
}

func (v value) isNumber() bool {
	return v.valueType == valNumber
}

func (v value) isString() bool {
	return v.valueType == valString
}

func boolValue(b bool) value {
	return value{valBool, b}
}

func nilValue() value {
	return value{valNil, 0}
}

func numberValue(n float64) value {
	return value{valNumber, n}
}

func stringValue(s string) value {
	return value{valString, s}
}

type valueArray []value

func (v *valueArray) write(b value) {
	*v = append(*v, b)
}

func (v value) isFalsey() bool {
	switch v.valueType {
	case valBool:
		return !v.value.(bool)
	case valNil:
		return true
	case valNumber:
		return v.value.(float64) == 0
	case valString:
		return v.value.(string) == ""
	default:
		return true
	}
}

func (v value) equals(other value) bool {
	return v.value == other.value
}

func (v value) String() string {
	switch v.valueType {
	case valBool:
		return fmt.Sprintf("\x1b[38;5;4m%v\x1b[0m", v.value)
	case valNil:
		return "\x1b[38;5;8m<nil>\x1b[0m"
	case valNumber:
		return fmt.Sprintf("\x1b[38;5;3m%v\x1b[0m", v.value)
	case valString:
		return fmt.Sprintf("\x1b[38;5;2m'%s'\x1b[0m", v.value)
	default:
		return fmt.Sprintf("%v", v.value)
	}
}
