package lox

import (
	"fmt"
	"slices"
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
		decl := p.declaration()
		if d, ok := decl.(*stmtExpr); ok && d.initializer == nil {
			break
		}
		p.program = append(p.program, decl)
	}
	return p.program, append(p.scanErrors, p.parseErrors...)
}

func (p *parser) expression() expression {
	return p.assignment()
}

func (p *parser) declaration() stmt {
	if p.match(Class) {
		return p.classDeclaration()
	}
	if p.match(Fun) {
		return p.function("function")
	}
	if p.match(Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *parser) classDeclaration() stmt {
	name := p.consume(Identifier, "Expected class name.")
	var methods []*stmtFun
	p.consume(LeftBrace, "Expected '{' before class body.")
	for !p.check(RightBrace) && !p.isAtEnd() {
		methods = append(methods, p.function("method").(*stmtFun))
	}
	p.consume(RightBrace, "Expected '}' after class body.")
	return &stmtClass{name, methods}
}

func (p *parser) function(kind string) stmt {
	name := p.consume(Identifier, "Expected "+kind+" name.")
	p.consume(LeftParen, "Expected '(' after "+kind+" name.")
	params := []token{}
	getParam := func() {
		if len(params) >= 255 {
			err := newError("Can't have more than 255 parameters.", p.peek().line)
			p.parseErrors = append(p.parseErrors, err)
		}
		params = append(params, p.consume(Identifier, "Expected parameter name."))
	}
	if !p.check(RightParen) {
		for getParam(); p.match(Comma); {
			getParam()
		}
	}
	p.consume(RightParen, "Expect ')' after parameters.")
	p.consume(LeftBrace, "Expect '{' before "+kind+" body.")
	body := p.blockStmt()
	return &stmtFun{name, body, params}
}

func (p *parser) varDeclaration() stmt {
	name := p.consume(Identifier, "Expected variable name.")
	var initializer expression
	if p.match(Equal) {
		initializer = p.expression()
	}
	p.consume(Semicolon, "Expected ';' after variable declaration.")
	return &stmtVar{initializer, name}
}

func (p *parser) statement() stmt {
	if p.match(For) {
		return p.forStmt()
	}
	if p.match(If) {
		return p.ifStmt()
	}
	if p.match(Return) {
		return p.returnStmt()
	}
	if p.match(While) {
		return p.whileStmt()
	}
	if p.match(LeftBrace) {
		return p.blockStmt()
	}
	expr := p.expression()
	p.consume(Semicolon, "Expected ';' after expression.")
	return &stmtExpr{expr}
}

func (p *parser) forStmt() stmt {
	p.consume(LeftParen, "Expect '(' after 'for'.")
	var initializer stmt
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Var) {
		initializer = p.varDeclaration()
	} else {
		initializer = &stmtExpr{p.expression()}
	}
	var condition expression
	if !p.check(Semicolon) {
		condition = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after loop condition.")
	var increment expression
	if !p.check(RightParen) {
		increment = p.expression()
	}
	p.consume(RightParen, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		body = &stmtBlock{[]stmt{body, &stmtExpr{increment}}}
	}
	if condition == nil {
		exprTrue := &exp{nil, nil, token{True, "true", "true", p.peek().line}}
		condition = &expressionLiteral{exprTrue, true}
	}
	body = &stmtWhile{condition, body}
	if initializer != nil {
		body = &stmtBlock{[]stmt{initializer, body}}
	}
	return body
}

func (p *parser) ifStmt() stmt {
	p.consume(LeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch stmt
	if p.match(Else) {
		elseBranch = p.statement()
	}
	return &stmtIf{condition, thenBranch, elseBranch}
}

func (p *parser) returnStmt() stmt {
	var val expression
	if !p.check(Semicolon) {
		val = p.expression()
	}
	p.consume(Semicolon, "Expected ';' after return value.")
	return &stmtReturn{val, p.previous()}
}

func (p *parser) whileStmt() stmt {
	p.consume(LeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after condition.")
	body := p.statement()
	return &stmtWhile{condition, body}
}

func (p *parser) blockStmt() stmt {
	stmts := []stmt{}
	for !p.check(RightBrace) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RightBrace, "Expected '}' after block.")
	return &stmtBlock{stmts}
}

func (p *parser) assignment() expression {
	expr := p.or()
	if p.match(Equal) {
		operator := p.previous()
		value := p.assignment()
		switch expr := expr.(type) {
		case *expressionVar:
			exp := &exp{expr, value, operator}
			return &expressionAssignment{exp}
		case *expressionGet:
			return &expressionSet{expr.expression, value, expr.name}
		case *expressionIndex:
			return &expressionSetIndex{expr.expression, expr.index, value}

		}
		err := newError("Invalid assignment target.", p.peek().line)
		p.parseErrors = append(p.parseErrors, err)
	}
	return expr
}

func (p *parser) or() expression {
	expr := p.and()
	for p.match(Or) {
		operator := p.previous()
		right := p.and()
		expr = &expressionLogical{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) and() expression {
	expr := p.equality()
	for p.match(And) {
		operator := p.previous()
		right := p.equality()
		expr = &expressionLogical{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) equality() expression {
	expr := p.comparison()
	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		return &expressionEquality{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) comparison() expression {
	expr := p.term()
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = &expressionComparison{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) term() expression {
	expr := p.factor()
	for p.match(Plus) {
		operator := p.previous()
		right := p.factor()
		expr = &expressionTerm{&exp{expr, right, operator}}
	}
	for p.match(Minus) {
		operator := p.previous()
		right := p.factor()
		expr = &expressionTerm{&exp{expr, right, operator}}
	}

	return expr
}

func (p *parser) factor() expression {
	expr := p.unary()
	for p.match(Star) {
		operator := p.previous()
		right := p.unary()
		expr = &expressionFactor{&exp{expr, right, operator}}
	}
	for p.match(Slash) {
		operator := p.previous()
		right := p.unary()
		expr = &expressionFactor{&exp{expr, right, operator}}
	}
	return expr
}

func (p *parser) unary() expression {
	if p.match(Bang) {
		operator := p.previous()
		right := p.unary()
		return &expressionUnary{&exp{nil, right, operator}}
	}
	if p.match(Minus) {
		operator := p.previous()
		right := p.unary()
		return &expressionUnary{&exp{nil, right, operator}}
	}

	return p.call()
}

func (p *parser) call() expression {
	expr := p.primary()
	for {
		if p.match(LeftParen) {
			expr = p.finishCall(expr)
		} else if p.match(Dot) {
			name := p.consume(Identifier, "Expect property name after '.'.")
			expr = &expressionGet{expr, name}
		} else if p.match(LeftBracket) {
			index := p.expression()
			p.consume(RightBracket, "Expect ']' after index.")
			expr = &expressionIndex{expr, index}
		} else {
			break
		}
	}
	return expr
}

func (p *parser) finishCall(callee expression) expression {
	args := []expression{}
	if !p.check(RightParen) {
		for {
			if len(args) >= 255 {
				err := newError("Can't have more than 255 arguments.", p.peek().line)
				p.parseErrors = append(p.parseErrors, err)
			}
			args = append(args, p.expression())
			if !p.match(Comma) {
				break
			}
		}
	}
	p.consume(RightParen, "Expect ')' after arguments.")
	return &expressionCall{callee, args}
}

func (p *parser) primary() expression {
	if p.match(False) {
		return &expressionLiteral{&exp{nil, nil, p.previous()}, false}
	}
	if p.match(True) {
		return &expressionLiteral{&exp{nil, nil, p.previous()}, true}
	}
	if p.match(Nil) {
		return &expressionLiteral{&exp{nil, nil, p.previous()}, nil}
	}
	if p.match(Number) {
		val := p.getFloatFromToken(p.previous().literal)
		return &expressionLiteral{&exp{nil, nil, p.previous()}, val}
	}
	if p.match(String) {
		val := p.previous().literal
		return &expressionLiteral{&exp{nil, nil, p.previous()}, val}
	}
	if p.match(LeftBracket) {
		items := []expression{}
		for !p.check(RightBracket) && !p.isAtEnd() {
			item := p.expression()
			if item == nil {
				break
			}
			items = append(items, item)
			if p.peek().tokenType != RightBracket {
				p.consume(Comma, "',' expected between array elements.")
			}
		}
		p.consume(RightBracket, "']' expected after array elements.")
		expr := &expressionLiteral{&exp{nil, nil, p.previous()}, items}
		return expr
	}
	if p.match(Identifier) {
		return &expressionVar{&exp{nil, nil, p.previous()}}
	}
	if p.match(This) {
		return &expressionThis{&exp{nil, nil, p.previous()}}
	}
	if p.match(LeftParen) {
		expr := &expressionGroup{p.expression()}
		p.consume(RightParen, "Unmatched parenthesis.")
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
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
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
		if p.previous().tokenType == Semicolon {
			return
		}
		switch p.peek().tokenType {
		case Class:
		case Fun:
		case Var:
		case For:
		case If:
		case While:
		case Return:
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
