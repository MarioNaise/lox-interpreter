package lox

import (
	"fmt"
	"strconv"
)

type parser struct {
	*scanner
	program     []stmtInterface
	parseErrors []loxError
	current     int
}

func newParser(str string) *parser {
	return &parser{scanner: newScanner(str)}
}

func (p *parser) parse() ([]stmtInterface, []loxError) {
	p.tokenize()
	for !p.isAtEnd() {
		p.program = append(p.program, p.declaration())
	}
	return p.program, append(p.scanErrors, p.parseErrors...)
}

func (p *parser) expression() exprInterface {
	return p.assignment()
}

func (p *parser) declaration() stmtInterface {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *parser) varDeclaration() stmtInterface {
	name := p.consume(IDENTIFIER, "Expected variable name.")
	var initializer exprInterface
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expected ';' after variable declaration.")
	return &stmtVar{&stmtExpr{initializer, name}}
}

func (p *parser) statement() stmtInterface {
	if p.match(FOR) {
		return p.forStmt()
	}
	if p.match(IF) {
		return p.ifStmt()
	}
	if p.match(PRINT) {
		return p.printStmt()
	}
	if p.match(WHILE) {
		return p.whileStmt()
	}
	if p.match(LEFT_BRACE) {
		return p.blockStmt()
	}
	expr := p.expression()
	p.consume(SEMICOLON, "Expected ';' after expression.")
	return &stmtExpr{initializer: expr}
}

func (p *parser) forStmt() stmtInterface {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer stmtInterface
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = &stmtExpr{initializer: p.expression()}
	}
	var condition exprInterface
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")
	var increment exprInterface
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		body = &stmtBlock{statements: []stmtInterface{body, &stmtExpr{initializer: increment}}}
	}
	if condition == nil {
		exprTrue := &expression{nil, nil, token{tokenType: TRUE, lexeme: "true", literal: "true", line: p.peek().line}}
		condition = &expressionLiteral{exprTrue, true}
	}
	body = &stmtWhile{condition: condition, body: body}
	if initializer != nil {
		body = &stmtBlock{statements: []stmtInterface{initializer, body}}
	}
	return body
}

func (p *parser) ifStmt() stmtInterface {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch stmtInterface
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return &stmtIf{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func (p *parser) printStmt() stmtInterface {
	value := p.expression()
	p.consume(SEMICOLON, "Expected ';' after value.")
	return &stmtPrint{&stmtExpr{initializer: value}}
}

func (p *parser) whileStmt() stmtInterface {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return &stmtWhile{condition: condition, body: body}
}

func (p *parser) blockStmt() stmtInterface {
	stmts := []stmtInterface{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expected '}' after block.")
	return &stmtBlock{statements: stmts}
}

func (p *parser) assignment() exprInterface {
	expr := p.or()
	if p.match(EQUAL) {
		operator := p.previous()
		value := p.assignment()
		switch expr := expr.(type) {
		case *expressionVar:
			exp := &expression{expr, value, operator}
			return &expressionAssignment{exp}
		}
		err := newError("Invalid assignment target.", p.peek().line)
		p.parseErrors = append(p.parseErrors, err)
	}
	return expr
}

func (p *parser) or() exprInterface {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &expressionLogical{&expression{expr, right, operator}}
	}
	return expr
}

func (p *parser) and() exprInterface {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &expressionLogical{&expression{expr, right, operator}}
	}
	return expr
}

func (p *parser) equality() exprInterface {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		return &expressionEquality{&expression{expr, right, operator}}
	}
	return expr
}

func (p *parser) comparison() exprInterface {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &expressionComparison{&expression{expr, right, operator}}
	}
	return expr
}

func (p *parser) term() exprInterface {
	expr := p.factor()
	for p.match(PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &expressionTerm{&expression{expr, right, operator}}
	}
	for p.match(MINUS) {
		operator := p.previous()
		right := p.factor()
		expr = &expressionTerm{&expression{expr, right, operator}}
	}

	return expr
}

func (p *parser) factor() exprInterface {
	expr := p.unary()
	for p.match(STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &expressionFactor{&expression{expr, right, operator}}
	}
	for p.match(SLASH) {
		operator := p.previous()
		right := p.unary()
		expr = &expressionFactor{&expression{expr, right, operator}}
	}
	return expr
}

func (p *parser) unary() exprInterface {
	if p.match(BANG) {
		operator := p.previous()
		right := p.unary()
		return &expressionUnary{&expression{nil, right, operator}}
	}
	if p.match(MINUS) {
		operator := p.previous()
		right := p.unary()
		return &expressionUnary{&expression{nil, right, operator}}
	}

	return p.primary()
}

func (p *parser) primary() exprInterface {
	if p.match(FALSE) {
		return &expressionLiteral{&expression{nil, nil, p.previous()}, false}
	}
	if p.match(TRUE) {
		return &expressionLiteral{&expression{nil, nil, p.previous()}, true}
	}
	if p.match(NIL) {
		return &expressionLiteral{&expression{nil, nil, p.previous()}, nil}
	}
	if p.match(NUMBER) {
		val := p.getFloatFromToken(p.previous().literal)
		return &expressionLiteral{&expression{nil, nil, p.previous()}, val}
	}
	if p.match(STRING) {
		val := p.previous().literal
		return &expressionLiteral{&expression{nil, nil, p.previous()}, val}
	}
	if p.match(IDENTIFIER) {
		return &expressionVar{&expression{nil, nil, p.previous()}}
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
	p.parseErrors = append(p.parseErrors, newError(err, p.previous().line))
	return token{}
}

func (p *parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}
		switch p.peek().tokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		p.advance()
	}
}

func (p *parser) getFloatFromToken(str string) float64 {
	valDouble, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(fmt.Sprintf("[line %d] Couldn't parse float64 -*%s*-", p.line, str))
	}
	return valDouble
}
