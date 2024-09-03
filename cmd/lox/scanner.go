package lox

import (
	"fmt"
	"regexp"
	"strings"
)

type regexRule struct {
	handler func(string)
	regex   string
}

type scanner struct {
	specCharTokenTypes map[string]string
	buffer             string
	tokens             []token
	scanErrors         []loxError
	regexRules         []regexRule
	keywords           []string
	specialChars       []string
	current            int
	line               int
}

func (s *scanner) tokenize() ([]token, []loxError) {
	s.tokens, s.scanErrors = []token{}, []loxError{}
	s.current = 0
	s.line = 1
outer:
	for s.current < len(s.buffer) {
		found := false
		for _, rule := range s.regexRules {
			r := regexp.MustCompile(`^` + rule.regex)
			loc := r.FindIndex([]byte(s.buffer[s.current:]))
			if len(loc) == 0 {
				continue
			}
			val := s.buffer[s.current : s.current+loc[1]]
			found = true
			rule.handler(val)
			s.current += len(val)
			continue outer
		}
		if !found {
			next := s.buffer[s.current : s.current+1]
			if next == `"` {
				s.scanErrors = append(s.scanErrors, newError("Unterminated string.", s.line))
				s.current++
				break outer
			} else {
				s.scanErrors = append(s.scanErrors, newError(fmt.Sprintf("Unexpected character: %s", next), s.line))
				s.current++
			}
		}
	}
	s.tokens = append(s.tokens, newToken(EOF, "", NULL, s.line))
	return s.tokens, s.scanErrors
}

func (s *scanner) defaultHandler(val string) {
	s.tokens = append(s.tokens, newToken(strings.ToUpper(val), val, NULL, s.line))
}

func (s *scanner) whitespaceHandler(val string) {
	if val == "\n" {
		s.line++
	}
}

func (s *scanner) specialCharHandler(val string) {
	s.tokens = append(s.tokens, newToken(s.specCharTokenTypes[val], val, NULL, s.line))
}

func (s *scanner) stringHandler(val string) {
	linesSkipped := len(regexp.MustCompile(`\n`).FindAllString(val, len(val)))
	s.line += linesSkipped
	s.tokens = append(s.tokens, newToken(STRING, val, val[1:len(val)-1], s.line))
}

func (s *scanner) identifierHandler(val string) {
	for _, keyword := range s.keywords {
		if strings.ToLower(keyword) == val {
			s.defaultHandler(val)
			return
		}
	}
	s.tokens = append(s.tokens, newToken(IDENTIFIER, val, NULL, s.line))
}

func (s *scanner) numberHandler(val string) {
	addComma := regexp.MustCompile(`(^\d+$)`)
	literal := addComma.ReplaceAllString(string(val), "$1.0")
	cutZeros := regexp.MustCompile(`([\d])0*$`)
	literal = cutZeros.ReplaceAllString(string(literal), "$1")
	s.tokens = append(s.tokens, newToken(NUMBER, val, literal, s.line))
}

func newScanner(str string) *scanner {
	l := &scanner{buffer: str}

	regexRules := []regexRule{
		{regex: "//.*", handler: func(_ string) {}},
		{regex: `\s`, handler: l.whitespaceHandler},
		{regex: `"[^"]*"`, handler: l.stringHandler},
		{regex: "[a-zA-Z_][a-zA-Z0-9_]*", handler: l.identifierHandler},
		{regex: `\d+(\.\d+)?`, handler: l.numberHandler},
	}

	l.keywords = []string{AND, CLASS, ELSE, FALSE, FOR, FUN, IF, NIL, OR, RETURN, SUPER, THIS, TRUE, VAR, WHILE}
	for _, keyword := range l.keywords {
		regexRules = append(regexRules, regexRule{regex: strings.ToLower(keyword), handler: l.defaultHandler})
	}

	l.specialChars = []string{`\!=`, `==`, `>=`, `<=`, `>`, `<`, `\!`, `=`, `;`, `\(`, `\)`, `{`, `}`, `\*`, `\.`, `,`, `\+`, `-`, `/`}
	for _, special := range l.specialChars {
		regexRules = append(regexRules, regexRule{regex: special, handler: l.specialCharHandler})
	}

	l.regexRules = regexRules
	l.specCharTokenTypes = map[string]string{
		"==": EQUAL_EQUAL,
		"!=": BANG_EQUAL,
		">=": GREATER_EQUAL,
		"<=": LESS_EQUAL,
		">":  GREATER,
		"<":  LESS,
		"!":  BANG,
		"=":  EQUAL,
		";":  SEMICOLON,
		"(":  LEFT_PAREN,
		")":  RIGHT_PAREN,
		"{":  LEFT_BRACE,
		"}":  RIGHT_BRACE,
		"*":  STAR,
		".":  DOT,
		",":  COMMA,
		"+":  PLUS,
		"-":  MINUS,
		"/":  SLASH,
	}

	return l
}
