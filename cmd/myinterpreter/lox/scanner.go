package lox

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type scanner struct {
	specaCharTokenTypes map[string]string
	commentString       string
	buffer              []rune
	specChars           []string
	keywords            []string
	tokens              []token
	errors              []loxError
}

// reads a list of tokens and possible errors from a given string
func (s *scanner) tokenize(str string) {
	defer func() {
		s.tokens = append(s.tokens, newToken(EOF, "", NULL))
	}()
	index := 0
	line := 1
	s.buffer = []rune(str)
	if len(s.buffer) == 0 {
		return
	}
	for index < len(s.buffer) {
		switch c := s.buffer[index]; {
		case s.isComment(index):
			index = s.handleComment(index)
		case isWhitespace(c):
			if c == '\n' {
				line++
			}
			index++
		case s.isSpecialChar(c):
			index = s.handleSpecialChar(index)
		case c == '"':
			index = s.handleString(index, line)
		case s.isKeywordStartLetter(c):
			index = s.handleKeywords(index)
		case unicode.IsLetter(c) || c == '_':
			index = s.handleIdentifier(index)
		case unicode.IsDigit(c):
			index = s.handleNumbers(index)
		default:
			s.errors = append(s.errors, newError("Unexpected character: "+string(c), line))
			index++
		}
	}
}

// returns true if char at index is the begin of a comment in buffer
func (s *scanner) isComment(index int) bool {
	runeLengthOfComment := len([]rune(s.commentString))
	if index+runeLengthOfComment-1 >= len(s.buffer) {
		return false
	}
	return string(s.buffer[index:index+runeLengthOfComment]) == s.commentString
}

// returns true if a character is a special character
func (s *scanner) isSpecialChar(char rune) bool {
	for _, specChar := range s.specChars {
		if char == []rune(specChar)[0] {
			return true
		}
	}
	return false
}

// returns true if any keywords start with the given letter
func (s *scanner) isKeywordStartLetter(c rune) bool {
	for _, keyW := range s.keywords {
		if c == []rune(keyW)[0] {
			return true
		}
	}
	return false
}

// returns true if a character is a whitespace character
func isWhitespace(c rune) bool {
	return unicode.IsSpace(c) || c == '\t' || c == '\n' || c == '\r'
}

// returns the index of the end of a comment
// returns the last index of the given []rune if no newline character is found
func (s *scanner) handleComment(index int) int {
	for i, char := range s.buffer[index:] {
		if char == '\n' {
			return index + i
		}
	}
	return len(s.buffer)
}

// handles special characters at given index in buffer
// returns the new index in buffer
func (s *scanner) handleSpecialChar(index int) int {
	for _, specChar := range s.specChars {
		if index+len([]rune(specChar)) > len(s.buffer) {
			continue
		}
		if string(s.buffer[index:index+len([]rune(specChar))]) == specChar {
			token := newToken(s.specaCharTokenTypes[specChar], string(specChar), NULL)
			s.tokens = append(s.tokens, token)
			return index + len([]rune(specChar))
		}
	}
	return index + 1
}

// checks for a keyword token at given index in buffer
// if not found, it gets handled as an identifier
// returns the new index in buffer
func (s *scanner) handleKeywords(index int) int {
	for _, keyW := range s.keywords {
		if index+len([]rune(keyW)) > len(s.buffer) {
			continue
		}
		if string(s.buffer[index:index+len([]rune(keyW))]) == keyW {
			token := newToken(strings.ToUpper(string(keyW)), string(keyW), NULL)
			s.tokens = append(s.tokens, token)
			return index + len([]rune(keyW))
		}
	}
	return s.handleIdentifier(index)
}

// handles identifier tokens at given index in buffer
// ends when it finds a non-letter, non-underscore char
// returns the new index in buffer
func (s *scanner) handleIdentifier(index int) int {
	var identifier []rune
	for _, char := range s.buffer[index:] {
		identifier = append(identifier, char)
		index++
		if index >= len(s.buffer) {
			break
		}
		c := s.buffer[index]
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
			break
		}
	}
	token := newToken(IDENTIFIER, string(identifier), NULL)
	s.tokens = append(s.tokens, token)
	return index
}

// handles string tokens at given index in buffer
// ends when it finds a non-escaped double quote
// returns the new index in buffer
func (s *scanner) handleString(index int, line int) int {
	var value []rune
	for _, char := range s.buffer[index:] {
		value = append(value, char)
		index++
		if index >= len(s.buffer) {
			break
		}
		c := s.buffer[index]
		if c == '"' && s.buffer[index-1] != '\\' {
			value = append(value, s.buffer[index])
			token := newToken(STRING, string(value), string(value[1:len(value)-1]))
			s.tokens = append(s.tokens, token)
			return index + 1
		}
	}
	s.errors = append(s.errors, newError("Unterminated string.", line))
	return len(s.buffer)
}

// handles number tokens at given index in buffer
// ends with the next non-digit character, except for '.'
func (s *scanner) handleNumbers(index int) int {
	var dotCount int
	var numberValue []rune
	for _, num := range s.buffer[index:] {
		numberValue = append(numberValue, num)
		index++
		if index >= len(s.buffer) {
			break
		}
		c := s.buffer[index]
		if c == '.' {
			dotCount++
			if dotCount > 1 {
				break
			}
			if index+1 >= len(s.buffer) || !unicode.IsDigit(s.buffer[index+1]) {
				break
			}
		}
		if !unicode.IsDigit(c) && c != '.' {
			break
		}
	}
	token := newToken(NUMBER, string(numberValue), func() string {
		if unicode.IsDigit(numberValue[len(numberValue)-1]) && !strings.Contains(string(numberValue), ".") {
			return fmt.Sprint(string(numberValue), ".0")
		}
		reg := regexp.MustCompile(`([\d])0*$`)
		return reg.ReplaceAllString(string(numberValue), "$1")
	}())

	s.tokens = append(s.tokens, token)
	return index
}

func newLoxScanner() scanner {
	return scanner{
		specaCharTokenTypes: map[string]string{
			"==": "EQUAL_EQUAL",
			"!=": "BANG_EQUAL",
			">=": "GREATER_EQUAL",
			"<=": "LESS_EQUAL",
			">":  "GREATER",
			"<":  "LESS",
			"!":  "BANG",
			"=":  "EQUAL",
			";":  "SEMICOLON",
			"(":  "LEFT_PAREN",
			")":  "RIGHT_PAREN",
			"{":  "LEFT_BRACE",
			"}":  "RIGHT_BRACE",
			"*":  "STAR",
			".":  "DOT",
			",":  "COMMA",
			"+":  "PLUS",
			"-":  "MINUS",
			"/":  "SLASH",
		},
		commentString: "//",
		specChars:     []string{"!=", "==", ">=", "<=", ">", "<", "!", "=", ";", "(", ")", "{", "}", "*", ".", ",", "+", "-", "/"},
		keywords:      []string{"and", "class", "else", "false", "for", "fun", "if", "nil", "or", "print", "return", "super", "this", "true", "var", "while"},
	}
}
