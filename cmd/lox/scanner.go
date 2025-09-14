package lox

import (
	"fmt"
	"strings"
)

var keywords = map[string]string{
	strings.ToLower(AND):    AND,
	strings.ToLower(CLASS):  CLASS,
	strings.ToLower(ELSE):   ELSE,
	strings.ToLower(FALSE):  FALSE,
	strings.ToLower(FOR):    FOR,
	strings.ToLower(FUN):    FUN,
	strings.ToLower(IF):     IF,
	strings.ToLower(NIL):    NIL,
	strings.ToLower(OR):     OR,
	strings.ToLower(RETURN): RETURN,
	strings.ToLower(THIS):   THIS,
	strings.ToLower(TRUE):   TRUE,
	strings.ToLower(VAR):    VAR,
	strings.ToLower(WHILE):  WHILE,
}

type scanner struct {
	source     []rune
	tokens     []token
	scanErrors []loxError
	start      int
	current    int
	line       int
}

func newScanner(str string) *scanner {
	s := &scanner{source: []rune(str), line: 1}
	return s
}

func (s *scanner) tokenize() ([]token, []loxError) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, newToken(EOF, string('\000'), NONE, s.line))
	return s.tokens, s.scanErrors
}

func (s *scanner) scanToken() {
	c := s.advance()
	switch c {
	case '[':
		s.addToken(LEFT_BRACKET, NONE)
	case ']':
		s.addToken(RIGHT_BRACKET, NONE)
	case '(':
		s.addToken(LEFT_PAREN, NONE)
	case ')':
		s.addToken(RIGHT_PAREN, NONE)
	case '{':
		s.addToken(LEFT_BRACE, NONE)
	case '}':
		s.addToken(RIGHT_BRACE, NONE)
	case ',':
		s.addToken(COMMA, NONE)
	case '.':
		s.addToken(DOT, NONE)
	case '-':
		s.addToken(MINUS, NONE)
	case '+':
		s.addToken(PLUS, NONE)
	case ';':
		s.addToken(SEMICOLON, NONE)
	case '*':
		s.addToken(STAR, NONE)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, NONE)
		} else {
			s.addToken(BANG, NONE)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, NONE)
		} else {
			s.addToken(EQUAL, NONE)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, NONE)
		} else {
			s.addToken(LESS, NONE)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, NONE)
		} else {
			s.addToken(GREATER, NONE)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, NONE)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.scanErrors = append(s.scanErrors, newError(fmt.Sprintf("Unexpected character: %c", c), s.line))
		}
	}
}

func (s *scanner) addToken(tokenType string, literal string) {
	s.tokens = append(s.tokens, newToken(tokenType, string(s.source[s.start:s.current]), literal, s.line))
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) advance() rune {
	s.current++
	return s.source[s.current-1]
}

func (s *scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.scanErrors = append(s.scanErrors, newError("Unterminated string.", s.line))
		return
	}

	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, string(value))
}

func (s *scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	s.addToken(NUMBER, string(s.source[s.start:s.current]))
}

func (s *scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := string(s.source[s.start:s.current])
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType, NONE)
}

func (s *scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
