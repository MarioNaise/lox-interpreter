package vm

import (
	"fmt"
	"os"
	"strconv"
)

type precedence int

type parser struct {
	current   token
	previous  token
	hadError  bool
	panicMode bool
}

const (
	precNone precedence = iota
	precAssignment
	precOr
	precAnd
	precEquality
	precComparison
	precTerm
	precFactor
	precUnary
	precendenceCall
	precPrimary
)

type parseRule struct {
	prefix     parseFn
	infix      parseFn
	precedence precedence
}

type parseFn func()

var (
	p              parser
	s              scanner
	compilingChunk = newChunk()
	rules          = []parseRule(nil)
)

func compile(src string, c *chunk) bool {
	s.init(src)
	compilingChunk = c
	p.hadError = false
	p.panicMode = false

	advance()
	expression()
	consume(tokenEOF, "Expect end of expression.")
	endCompiler()
	return !p.hadError
}

func errorAt(token token, message string) {
	if p.panicMode {
		return
	}
	p.panicMode = true

	fmt.Fprintf(os.Stderr, "[line %d] Error", token.line)

	switch token.tokenType {
	case tokenEOF:
		fmt.Fprintf(os.Stderr, " at end")
	case tokenError:
	default:
		fmt.Fprintf(os.Stderr, " at '%s'. ", token.lexeme)
	}
	fmt.Fprintln(os.Stderr, message)
	p.hadError = true
}

func error(message string) {
	errorAt(p.previous, message)
}

func errorAtCurrent(message string) {
	errorAt(p.current, message)
}

func currentChunk() *chunk {
	return compilingChunk
}

func advance() {
	p.previous = p.current
	for {
		p.current = s.scanToken()
		if p.current.tokenType != tokenError {
			break
		}
		errorAtCurrent(p.current.lexeme)
	}
}

func consume(tokenType tokenType, message string) {
	if p.current.tokenType == tokenType {
		advance()
		return
	}
	errorAtCurrent(message)
}

func emitByte(b byte) {
	currentChunk().write(b, p.previous.line)
}

func emitBytes(b1, b2 byte) {
	emitByte(b1)
	emitByte(b2)
}

func emitReturn() {
	emitByte(byte(opReturn))
}

func makeConstant(value value) byte {
	constant := currentChunk().addConstant(value)
	if constant > 255 {
		error("Too many constants in one chunk.")
		return 0
	}
	return byte(constant)
}

func emitConstant(value value) {
	emitBytes(byte(opConstant), makeConstant(value))
}

func endCompiler() {
	emitReturn()

	if debug {
		if !p.hadError {
			currentChunk().disassemble("code")
		}
	}
}

func getRule(tokenType tokenType) parseRule {
	if rules == nil {
		rules = []parseRule{
			{grouping, nil, precNone},     // tokenLeftParen
			{nil, nil, precNone},          // tokenRightParen
			{nil, nil, precNone},          // tokenLeftBrace
			{nil, nil, precNone},          // tokenRightBrace
			{nil, nil, precNone},          // tokenComma
			{nil, nil, precNone},          // tokenDot
			{unary, binary, precTerm},     // tokenMinus
			{nil, binary, precTerm},       // tokenPlus
			{nil, nil, precNone},          // tokenSemicolon
			{nil, binary, precFactor},     // tokenSlash
			{nil, binary, precFactor},     // tokenStar
			{unary, nil, precNone},        // tokenBang
			{nil, binary, precEquality},   // tokenBangEqual
			{nil, nil, precNone},          // tokenEqual
			{nil, binary, precEquality},   // tokenEqualEqual
			{nil, binary, precComparison}, // tokenGreater
			{nil, binary, precComparison}, // tokenGreaterEqual
			{nil, binary, precComparison}, // tokenLess
			{nil, binary, precComparison}, // tokenLessEqual
			{nil, nil, precNone},          // tokenIdentifier
			{str, nil, precNone},          // tokenString
			{number, nil, precNone},       // tokenNumber
			{nil, nil, precNone},          // tokenAnd
			{nil, nil, precNone},          // tokenClass
			{nil, nil, precNone},          // tokenElse
			{literal, nil, precNone},      // tokenFalse
			{nil, nil, precNone},          // tokenFor
			{nil, nil, precNone},          // tokenFun
			{nil, nil, precNone},          // tokenIf
			{literal, nil, precNone},      // tokenNil
			{nil, nil, precNone},          // tokenOr
			{nil, nil, precNone},          // tokenPrint
			{nil, nil, precNone},          // tokenReturn
			{nil, nil, precNone},          // tokenSuper
			{nil, nil, precNone},          // tokenThis
			{literal, nil, precNone},      // tokenTrue
			{nil, nil, precNone},          // tokenVar
			{nil, nil, precNone},          // tokenWhile
			{nil, nil, precNone},          // tokenError
			{nil, nil, precNone},          // tokenEOF
		}
	}
	return rules[tokenType]
}

func binary() {
	opType := p.previous.tokenType
	parseRule := getRule(opType)
	parsePrecedence(parseRule.precedence + 1)

	switch opType {
	case tokenBangEqual:
		emitBytes(byte(opEqual), byte(opNot))
	case tokenEqualEqual:
		emitByte(byte(opEqual))
	case tokenGreater:
		emitByte(byte(opGreater))
	case tokenGreaterEqual:
		emitBytes(byte(opLess), byte(opNot))
	case tokenLess:
		emitByte(byte(opLess))
	case tokenLessEqual:
		emitBytes(byte(opGreater), byte(opNot))
	case tokenPlus:
		emitByte(byte(opAdd))
	case tokenMinus:
		emitByte(byte(opSubtract))
	case tokenStar:
		emitByte(byte(opMultiply))
	case tokenSlash:
		emitByte(byte(opDivide))
	default:
		panic("Invalid binary operator")
	}
}

func literal() {
	switch p.previous.tokenType {
	case tokenFalse:
		emitByte(byte(opFalse))
	case tokenNil:
		emitByte(byte(opNil))
	case tokenTrue:
		emitByte(byte(opTrue))
	default:
		panic("Invalid literal")
	}
}

func grouping() {
	expression()
	consume(tokenRightParen, "Expect ')' after expression.")
}

func number() {
	val, err := strconv.ParseFloat(p.previous.lexeme, 64)
	if err != nil {
		panic("Invalid number")
	}
	emitConstant(numberValue(val))
}

func str() {
	emitConstant(stringValue(p.previous.lexeme[1 : len(p.previous.lexeme)-1]))
}

func unary() {
	opType := p.previous.tokenType

	// Compile the operand.
	parsePrecedence(precUnary)

	// Emit the operator instruction.
	switch opType {
	case tokenBang:
		emitByte(byte(opNot))
	case tokenMinus:
		emitByte(byte(opNegate))
	default:
		panic("Invalid unary operator")
	}
}

func parsePrecedence(precedence precedence) {
	advance()
	prefixRule := getRule(p.previous.tokenType).prefix
	if prefixRule == nil {
		error("Expect expression.")
		return
	}

	prefixRule()

	for precedence <= getRule(p.current.tokenType).precedence {
		advance()
		infixRule := getRule(p.previous.tokenType).infix
		infixRule()
	}
}

func expression() {
	parsePrecedence(precAssignment)
}
