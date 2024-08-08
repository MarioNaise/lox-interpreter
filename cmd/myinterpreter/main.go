package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type (
	Token       string
	Error       string
	Keyword     string
	SpecialChar string
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	text := string(fileContents)
	tokens, errors := tokenize(text)

	for _, token := range tokens {
		fmt.Println(token)
	}
	if len(errors) > 0 {
		for _, error := range errors {
			fmt.Fprintf(os.Stderr, "%s\n", error)
		}
		os.Exit(65)
	}
	os.Exit(0)
}

// returns a list or tokens and errors from a string
func tokenize(s string) ([]Token, []Error) {
	line := 1
	var tokens []Token
	var errors []Error
	if s == "" {
		return append(tokens, "EOF  null"), errors
	}
	chars := []rune(s)
	var pos int
	for pos < len(chars) {
		switch c := chars[pos]; {
		case isComment(chars, pos):
			pos = handleComment(chars, pos)
		case isWhitespace(c):
			if c == '\n' {
				line++
			}
			pos++
		case isSpecialChar(c):
			pos, tokens = handleSpecialChar(chars, pos, tokens)
		case c == '"':
			pos, tokens, errors = handleString(chars, pos, line, tokens, errors)
		case isKeywordStartLetter(c):
			pos, tokens = handleKeywords(chars, pos, tokens)
		case unicode.IsLetter(c) || c == '_':
			pos, tokens = handleIdentifier(chars, pos, tokens)
		default:
			errors = append(errors, Error(fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, c)))
			pos++
		}
	}
	return append(tokens, "EOF  null"), errors
}

// returns the list of keywords
func getKeywords() [1]Keyword {
	return [1]Keyword{"var"}
}

// returns the list of special characters
func getSpecialChars() []SpecialChar {
	return []SpecialChar{"!=", "==", ">=", "<=", ">", "<", "!", "=", ";", "(", ")", "{", "}", "*", ".", ",", "+", "-", "/"}
}

// returns the mapping for special characters
func getCharMap() map[SpecialChar]string {
	return map[SpecialChar]string{
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
	}
}

// returns true if currenct + next character equals to "//"
func isComment(chars []rune, pos int) bool {
	if pos+1 >= len(chars) {
		return false
	}
	return chars[pos] == '/' && chars[pos+1] == '/'
}

// returns true if a character is a special character
func isSpecialChar(char rune) bool {
	for _, specialChar := range getSpecialChars() {
		if char == rune(specialChar[0]) {
			return true
		}
	}
	return false
}

// returns true if any keywords start with the given letter
func isKeywordStartLetter(c rune) bool {
	for _, keyword := range getKeywords() {
		if c == rune(keyword[0]) {
			return true
		}
	}
	return false
}

// returns true if a character is a whitespace character
func isWhitespace(c rune) bool {
	return unicode.IsSpace(c) || c == '\t' || c == '\n' || c == '\r'
}

// returns the position of the end of a comment
// returns the last position of the given []rune if no newline character is found
func handleComment(chars []rune, pos int) int {
	for i, char := range chars[pos:] {
		if char == '\n' {
			return pos + i
		}
	}
	return len(chars)
}

// handles special characters at given position in chars
// returns the new position in chars
func handleSpecialChar(chars []rune, pos int, tokens []Token) (int, []Token) {
	charMap := getCharMap()
	specialChars := getSpecialChars()
	for _, specialChar := range specialChars {
		if SpecialChar(chars[pos:pos+len(specialChar)]) == specialChar {
			token := Token(fmt.Sprintf("%s %s %s", charMap[specialChar], specialChar, "null"))
			return pos + len(specialChar), append(tokens, token)
		}
	}
	return pos + 1, tokens
}

// checks chars for a keyword token at given position in chars
// if not found, it gets handled as an identifier
// returns the new position in chars
func handleKeywords(chars []rune, pos int, tokens []Token) (int, []Token) {
	for _, keyword := range getKeywords() {
		if Keyword(chars[pos:pos+len(keyword)]) == keyword {
			token := Token(fmt.Sprintf("%s %s %s", strings.ToUpper(string(keyword)), keyword, "null"))
			return pos + len(keyword), append(tokens, token)
		}
	}
	return handleIdentifier(chars, pos, tokens)
}

// handles identifier tokens at given position in chars
// ends when it finds a non-letter, non-underscore char
// returns the new position in chars
func handleIdentifier(chars []rune, pos int, tokens []Token) (int, []Token) {
	var identifier []rune
	for {
		identifier = append(identifier, chars[pos])
		pos++
		c := chars[pos]
		if !unicode.IsLetter(c) && c != '_' {
			token := Token(fmt.Sprintf("IDENTIFIER %s %s", string(identifier), "null"))
			return pos + 1, append(tokens, token)
		}
	}
}

// handles string tokens at given position in chars
// ends when it finds a non-escaped double quote
// returns the new position in chars
func handleString(chars []rune, pos int, line int, tokens []Token, errors []Error) (int, []Token, []Error) {
	var value []rune
	for _, char := range chars[pos:] {
		value = append(value, char)
		pos++
		if pos >= len(chars) {
			break
		}
		c := chars[pos]
		if c == '"' && rune(chars[pos-1]) != '\\' {
			value = append(value, chars[pos])
			token := Token(fmt.Sprintf("STRING %s %s", string(value), string(value[1:len(value)-1])))
			return pos + 1, append(tokens, token), errors
		}
	}
	return len(chars), tokens, append(errors, Error(fmt.Sprintf("[line %d] Error: Unterminated string.", line)))
}
