package lox

import (
	"io"
)

type parser struct {
	*scanner
	expression  exprInterface
	parseErrors []loxError
	current     int
}

func (p *parser) parse() {
	p.tokenize()
	p.expression = p.equality()
}

func (p *parser) equality() exprInterface {
	expr = p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &expression{expr, right, operator}
	}
	return expr
}

func (p *parser) comparison() exprInterface {
	expr = p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &expression{expr, right, operator}
	}
	return expr
}

func (p *parser) term() exprInterface {
	expr = p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &expression{expr, right, operator}
	}
	return expr
}

func (p *parser) factor() exprInterface {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &expression{expr, right, operator}
	}
	return expr
}

func (p *parser) unary() exprInterface {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return &expression{nil, right, operator}
	}
	return p.primary()
}

func (p *parser) primary() exprInterface {
	if p.match(FALSE, TRUE, NIL, NUMBER, STRING) {
		return &expression{nil, nil, p.previous()}
	}
	if p.match(LEFT_PAREN) {
		expr := &groupExpression{p.equality()}
		p.consume(RIGHT_PAREN, "Unmatched parenthesis.")
		return expr
	}
	p.parseErrors = append(p.parseErrors, newError("at '"+p.peek().lexeme+"' - Expected expression.", p.peek().line))
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
	return p.current >= len(p.tokens)
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

func newParser(r io.Reader) *parser {
	return &parser{scanner: newScanner(r)}
}
