package lox

import (
	"fmt"
	"strconv"
)

type parser struct {
	*scanner
	program     []stmt
	parseErrors []loxError
	current     int
}

func newParser(str string) *parser {
	return &parser{scanner: newScanner(str)}
}

func (p *parser) parse() ([]stmt, []loxError) {
	p.tokenize()
	for !p.isAtEnd() {
		p.program = append(p.program, p.declaration())
	}
	return p.program, append(p.scanErrors, p.parseErrors...)
}

func (p *parser) expression() expression {
	return p.assignment()
}

func (p *parser) declaration() stmt {
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *parser) classDeclaration() stmt {
	name := p.consume(IDENTIFIER, "Expected class name.")
	var methods []stmtFun
	p.consume(LEFT_BRACE, "Expected '{' before class body.")
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, *p.function("method").(*stmtFun))
	}
	p.consume(RIGHT_BRACE, "Expected '}' after class body.")
	return &stmtClass{name, methods}
}

func (p *parser) function(kind string) stmt {
	name := p.consume(IDENTIFIER, "Expected "+kind+" name.")
	p.consume(LEFT_PAREN, "Expected '(' after "+kind+" name.")
	params := []token{}
	getParam := func() {
		if len(params) >= 255 {
			err := newError("Can't have more than 255 parameters.", p.peek().line)
			p.parseErrors = append(p.parseErrors, err)
		}
		params = append(params, p.consume(IDENTIFIER, "Expected parameter name."))
	}
	if !p.check(RIGHT_PAREN) {
		for getParam(); p.match(COMMA); {
			getParam()
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.blockStmt()
	return &stmtFun{name, body, params}
}

func (p *parser) varDeclaration() stmt {
	name := p.consume(IDENTIFIER, "Expected variable name.")
	var initializer expression
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expected ';' after variable declaration.")
	return &stmtVar{initializer, name}
}

func (p *parser) statement() stmt {
	if p.match(FOR) {
		return p.forStmt()
	}
	if p.match(IF) {
		return p.ifStmt()
	}
	if p.match(RETURN) {
		return p.returnStmt()
	}
	if p.match(WHILE) {
		return p.whileStmt()
	}
	if p.match(LEFT_BRACE) {
		return p.blockStmt()
	}
	expr := p.expression()
	p.consume(SEMICOLON, "Expected ';' after expression.")
	return &stmtExpr{expr}
}

func (p *parser) forStmt() stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = &stmtExpr{p.expression()}
	}
	var condition expression
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")
	var increment expression
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		body = &stmtBlock{[]stmt{body, &stmtExpr{increment}}}
	}
	if condition == nil {
		exprTrue := &exp{nil, nil, token{TRUE, "true", "true", p.peek().line}}
		condition = &expressionLiteral{exprTrue, true}
	}
	body = &stmtWhile{condition, body}
	if initializer != nil {
		body = &stmtBlock{[]stmt{initializer, body}}
	}
	return body
}

func (p *parser) ifStmt() stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return &stmtIf{condition, thenBranch, elseBranch}
}

func (p *parser) returnStmt() stmt {
	var val expression
	if !p.check(SEMICOLON) {
		val = p.expression()
	}
	p.consume(SEMICOLON, "Expected ';' after return value.")
	return &stmtReturn{val, p.previous()}
}

func (p *parser) whileStmt() stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return &stmtWhile{condition, body}
}

func (p *parser) blockStmt() stmt {
	stmts := []stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expected '}' after block.")
	return &stmtBlock{stmts}
}

func (p *parser) assignment() expression {
	expr := p.or()
	if p.match(EQUAL) {
		operator := p.previous()
		value := p.assignment()
		switch expr := expr.(type) {
		case *expressionVar:
			exp := &exp{expr, value, operator}
			return &expressionAssignment{exp}
		}
		err := newError("Invalid assignment target.", p.peek().line)
		p.parseErrors = append(p.parseErrors, err)
	}
	return expr
}

func (p *parser) or() expression {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &expressionLogical{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) and() expression {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &expressionLogical{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) equality() expression {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		return &expressionEquality{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) comparison() expression {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &expressionComparison{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) term() expression {
	expr := p.factor()
	for p.match(PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &expressionTerm{&exp{expr, right, operator}}
	}
	for p.match(MINUS) {
		operator := p.previous()
		right := p.factor()
		expr = &expressionTerm{&exp{expr, right, operator}}
	}

	return expr
}

func (p *parser) factor() expression {
	expr := p.unary()
	for p.match(STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &expressionFactor{&exp{expr, right, operator}}
	}
	for p.match(SLASH) {
		operator := p.previous()
		right := p.unary()
		expr = &expressionFactor{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) unary() expression {
	if p.match(BANG) {
		operator := p.previous()
		right := p.unary()
		return &expressionUnary{&exp{nil, right, operator}}
	}
	if p.match(MINUS) {
		operator := p.previous()
		right := p.unary()
		return &expressionUnary{&exp{nil, right, operator}}
	}

	return p.call()
}

func (p *parser) call() expression {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *parser) finishCall(callee expression) expression {
	args := []expression{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(args) >= 255 {
				err := newError("Can't have more than 255 arguments.", p.peek().line)
				p.parseErrors = append(p.parseErrors, err)
			}
			args = append(args, p.expression())
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	return &expressionCall{callee, args}
}

func (p *parser) primary() expression {
	if p.match(FALSE) {
		return &expressionLiteral{&exp{nil, nil, p.previous()}, false}
	}
	if p.match(TRUE) {
		return &expressionLiteral{&exp{nil, nil, p.previous()}, true}
	}
	if p.match(NIL) {
		return &expressionLiteral{&exp{nil, nil, p.previous()}, nil}
	}
	if p.match(NUMBER) {
		val := p.getFloatFromToken(p.previous().literal)
		return &expressionLiteral{&exp{nil, nil, p.previous()}, val}
	}
	if p.match(STRING) {
		val := p.previous().literal
		return &expressionLiteral{&exp{nil, nil, p.previous()}, val}
	}
	if p.match(IDENTIFIER) {
		return &expressionVar{&exp{nil, nil, p.previous()}}
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
