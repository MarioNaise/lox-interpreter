package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type regexRule struct {
	handler func(string)
	regex   string
}

type scanner struct {
	*bufio.Scanner
	specCharTokenTypes map[string]string
	tokens             []token
	scanErrors         []loxError
	regexRules         []regexRule
	keywords           []string
	specialChars       []string
	current            int
	line               int
}

func (s *scanner) tokenize() {
	for s.Scan() {
		s.current = 0
		s.line++
		s.ScanLine()

	}
	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s.tokens = append(s.tokens, newToken(EOF, "", NULL, s.line))
}

func (s *scanner) ScanLine() {
outer:
	for s.current < len(s.Text()) {
		valueFound := false
		for _, rule := range s.regexRules {
			reg := regexp.MustCompile(rule.regex)
			loc := reg.FindIndex([]byte(s.Text()[s.current:]))
			if len(loc) < 2 || loc[0] != 0 {
				continue
			}
			valueFound = true
			value := string(s.Text()[s.current : s.current+loc[1]])
			rule.handler(value)
			s.current += len(value)
			continue outer
		}
		if !valueFound {
			if s.Text()[s.current:s.current+1] == `"` {
				s.scanErrors = append(s.scanErrors, newError("Unterminated string.", s.line))
				s.current++
				break outer
			} else {
				s.scanErrors = append(s.scanErrors, newError("Unexpected character: "+string(s.Text()[s.current]), s.line))
				s.current++
			}
		}
	}
}

func (s *scanner) defaultHandler(val string) {
	s.tokens = append(s.tokens, newToken(strings.ToUpper(val), val, NULL, s.line))
}

func (s *scanner) specialCharHandler(val string) {
	s.tokens = append(s.tokens, newToken(s.specCharTokenTypes[val], val, NULL, s.line))
}

func (s *scanner) stringHandler(val string) {
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

func newScanner(r io.Reader) *scanner {
	l := &scanner{Scanner: bufio.NewScanner(r)}

	regexRules := []regexRule{
		{regex: "//.*", handler: func(_ string) {}},
		{regex: `\s+`, handler: func(_ string) {}},
		{regex: `"[^"]*"`, handler: l.stringHandler},
		{regex: "[a-zA-Z_][a-zA-Z0-9_]*", handler: l.identifierHandler},
		{regex: `\d+(\.\d+)?`, handler: l.numberHandler},
	}

	l.keywords = []string{AND, CLASS, ELSE, FALSE, FOR, FUN, IF, NIL, OR, PRINT, RETURN, SUPER, THIS, TRUE, VAR, WHILE}
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
