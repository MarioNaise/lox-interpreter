package lox

import (
	"fmt"
	"strings"
)

var keywords = map[string]string{
	strings.ToLower(And):    And,
	strings.ToLower(Class):  Class,
	strings.ToLower(Else):   Else,
	strings.ToLower(False):  False,
	strings.ToLower(For):    For,
	strings.ToLower(Fun):    Fun,
	strings.ToLower(If):     If,
	strings.ToLower(Nil):    Nil,
	strings.ToLower(Or):     Or,
	strings.ToLower(Return): Return,
	strings.ToLower(This):   This,
	strings.ToLower(True):   True,
	strings.ToLower(Var):    Var,
	strings.ToLower(While):  While,
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
	s.tokens = append(s.tokens, newToken(EOF, string('\000'), None, s.line))
	return s.tokens, s.scanErrors
}

func (s *scanner) scanToken() {
	c := s.advance()
	switch c {
	case '[':
		s.addToken(LeftBracket, None)
	case ']':
		s.addToken(RightBracket, None)
	case '(':
		s.addToken(LeftParen, None)
	case ')':
		s.addToken(RightParen, None)
	case '{':
		s.addToken(LeftBrace, None)
	case '}':
		s.addToken(RightBrace, None)
	case ',':
		s.addToken(Comma, None)
	case '.':
		s.addToken(Dot, None)
	case '-':
		s.addToken(Minus, None)
	case '+':
		s.addToken(Plus, None)
	case ';':
		s.addToken(Semicolon, None)
	case '*':
		s.addToken(Star, None)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual, None)
		} else {
			s.addToken(Bang, None)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual, None)
		} else {
			s.addToken(Equal, None)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual, None)
		} else {
			s.addToken(Less, None)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual, None)
		} else {
			s.addToken(Greater, None)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash, None)
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
	currentLine := s.line
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.scanErrors = append(s.scanErrors, newError("Unterminated string.", currentLine))
		return
	}

	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(String, string(value))
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
	s.addToken(Number, string(s.source[s.start:s.current]))
}

func (s *scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := string(s.source[s.start:s.current])
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = Identifier
	}
	s.addToken(tokenType, None)
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
