package lox

type expressionGroup struct {
	exprInterface
}

func (e expressionGroup) String() string {
	return "(group " + e.exprInterface.String() + ")"
}

func (e *expressionGroup) evaluate() string {
	return e.exprInterface.evaluate()
}
