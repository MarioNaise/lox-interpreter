package lox

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type parser struct {
	*scanner
	program     []stmtInterface
	parseErrors []loxError
	current     int
}

func newParser(r io.Reader) *parser {
	return &parser{scanner: newScanner(r)}
}

func (p *parser) parse() ([]stmtInterface, []loxError) {
	p.parseErrors = []loxError{}
	p.tokenize()
	for !p.isAtEnd() {
		p.program = append(p.program, p.statement())
	}
	return p.program, p.parseErrors
}

func (p *parser) statement() stmtInterface {
	if p.match(PRINT) {
		return p.printStmt()
	}
	expr := p.equality()
	p.consume(SEMICOLON, "Expected ';' after expression.")
	return &stmtExpr{expr}
}

func (p *parser) printStmt() stmtInterface {
	value := p.equality()
	p.consume(SEMICOLON, "Expected ';' after value.")
	return &stmtPrint{&stmtExpr{value}}
}

func (p *parser) equality() exprInterface {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		return &expressionEquality{&expression{expr, right, false, operator}}
	}
	return expr
}

func (p *parser) comparison() exprInterface {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &expressionComparison{&expression{expr, right, false, operator}}
	}
	return expr
}

func (p *parser) term() exprInterface {
	expr := p.factor()
	for p.match(PLUS) {
		operator := p.previous()
		right := p.factor()
		var val any
		if reflect.TypeOf(expr.value()).Name() == "string" ||
			reflect.TypeOf(right.value()).Name() == "string" {
			val = fmt.Sprintf("%v%v", expr.value(), right.value())
		} else {
			val = p.getFloatFromValue(expr) + p.getFloatFromValue(right)
		}
		expr = &expressionTerm{&expression{expr, right, val, operator}}
	}
	for p.match(MINUS) {
		operator := p.previous()
		right := p.factor()
		val := p.getFloatFromValue(expr) - p.getFloatFromValue(right)
		expr = &expressionTerm{&expression{expr, right, val, operator}}
	}

	return expr
}

func (p *parser) factor() exprInterface {
	expr := p.unary()
	for p.match(STAR) {
		operator := p.previous()
		right := p.unary()
		val := p.getFloatFromValue(expr) * p.getFloatFromValue(right)
		expr = &expressionFactor{&expression{expr, right, val, operator}}
	}
	for p.match(SLASH) {
		operator := p.previous()
		right := p.unary()
		val := p.getFloatFromValue(expr) / p.getFloatFromValue(right)
		expr = &expressionFactor{&expression{expr, right, val, operator}}
	}
	return expr
}

func (p *parser) unary() exprInterface {
	if p.match(BANG) {
		operator := p.previous()
		right := p.unary()
		val := p.isTruthy(right)
		return &expressionUnary{&expression{nil, right, val, operator}}
	}
	if p.match(MINUS) {
		operator := p.previous()
		right := p.unary()
		val := -p.getFloatFromValue(right)
		return &expressionUnary{&expression{nil, right, val, operator}}
	}

	return p.primary()
}

func (p *parser) primary() exprInterface {
	if p.match(FALSE) {
		return &expressionLiteral{&expression{nil, nil, false, p.previous()}}
	}
	if p.match(TRUE) {
		return &expressionLiteral{&expression{nil, nil, true, p.previous()}}
	}
	if p.match(NIL) {
		return &expressionLiteral{&expression{nil, nil, nil, p.previous()}}
	}
	if p.match(NUMBER) {
		val := p.getFloatFromToken(p.previous().literal)
		return &expressionLiteral{&expression{nil, nil, val, p.previous()}}
	}
	if p.match(STRING) {
		val := p.previous().literal
		return &expressionLiteral{&expression{nil, nil, val, p.previous()}}
	}
	if p.match(LEFT_PAREN) {
		expr := &expressionGroup{p.equality()}
		p.consume(RIGHT_PAREN, "Unmatched parenthesis.")
		return expr
	}
	err := newError("at '"+p.peek().lexeme+"' - Expected expression.", p.peek().line)
	p.parseErrors = append(p.parseErrors, err)
	return nil
}

/////////////////////
/// Helper methods///
/////////////////////

func (p *parser) advance() {
	if !p.isAtEnd() {
		p.current++
	}
}

func (p *parser) peek() token {
	return p.tokens[p.current]
}

func (p *parser) previous() token {
	return p.tokens[p.current-1]
}

func (p *parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *parser) match(types ...string) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(t string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == t
}

func (p *parser) consume(t string, err string) token {
	if p.peek().tokenType == t {
		p.advance()
		return p.previous()
	}
	p.parseErrors = append(p.parseErrors, newError(err, p.peek().line))
	return token{}
}

func (p *parser) getFloatFromToken(str string) float64 {
	valDouble, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(fmt.Sprintf("[line %d] Couldn't parse float64 -*%s*-", p.line, str))
	}
	return valDouble
}

func (p *parser) getFloatFromValue(e exprInterface) float64 {
	if reflect.TypeOf(e.value()).Name() != "float64" {
		if p.isTruthyByType(e) {
			return 1
		}
		if !p.isTruthyByType(e) {
			return 0
		}
	}
	return e.value().(float64)
}

func (p *parser) isTruthy(e exprInterface) bool {
	value := e.value()
	if value == nil {
		return false
	}
	switch e.tokenType() {
	case BANG:
		return !p.isTruthyByType(e)
	case STRING:
		return value != ""
	case NUMBER:
		return value.(float64) != 0
	case TRUE:
		return true
	case FALSE:
		return false
	case NIL:
		return false
	}
	return p.isTruthyByType(e)
}

func (p *parser) isTruthyByType(e exprInterface) bool {
	value := e.value()
	switch reflect.TypeOf(value).Name() {
	case "string":
		return value != ""
	case "float64":
		return value.(float64) != 0
	case "bool":
		return value.(bool)
	}
	return false
}
