package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Token string

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
	tokens := tokenize(text)
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func tokenize(s string) []Token {
	var tokens []Token
	if s == "" {
		return append(tokens, "EOF  null")
	}
	chars := []rune(s)
	var pos int
	for pos < len(chars) {
		switch c := chars[pos]; {
		case isWhitespace(c):
			pos++
		case isSpecialChar(c):
			pos, tokens = handleSpecialChar(c, pos, tokens)
		case c == '"':
			pos, tokens = handleString(chars, pos, tokens)
		case isKeywordStartLetter(c):
			pos, tokens = handleKeywords(chars, pos, tokens)
		case unicode.IsLetter(c) || c == '_':
			pos, tokens = handleIdentifier(chars, pos, tokens)
		default:
			pos++
		}
	}
	return append(tokens, "EOF  null")
}

// returns the list of special characters
func getSpecialChars() [6]rune {
	return [6]rune{'=', ';', '(', ')', '{', '}'}
}

func getCharMap() map[rune]string {
	return map[rune]string{
		'=': "EQUAL",
		';': "SEMICOLON",
		'(': "LEFT_PAREN",
		')': "RIGHT_PAREN",
		'{': "LEFT_BRACE",
		'}': "RIGHT_BRACE",
	}
}

// returns the list of keywords
func getKeywords() [1]string {
	return [1]string{"var"}
}

// returns true if a character is a special character
func isSpecialChar(char rune) bool {
	for _, specialChar := range getSpecialChars() {
		if char == specialChar {
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

// handles special characters at given position in chars
// returns the new position in chars
func handleSpecialChar(char rune, pos int, tokens []Token) (int, []Token) {
	charMap := getCharMap()
	token := Token(fmt.Sprintf("%s %s %s", charMap[char], string(char), "null"))
	return pos + 1, append(tokens, token)
}

// checks chars for a keyword token at given position in chars
// if not found, it gets handled as an identifier
// returns the new position in chars
func handleKeywords(chars []rune, pos int, tokens []Token) (int, []Token) {
	for _, keyword := range getKeywords() {
		if string(chars[pos:pos+len(keyword)]) == keyword {
			token := Token(fmt.Sprintf("%s %s %s", strings.ToUpper(keyword), keyword, "null"))
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
func handleString(chars []rune, pos int, tokens []Token) (int, []Token) {
	var value []rune
	for {
		value = append(value, chars[pos])
		pos++
		c := chars[pos]
		if c == '"' && rune(chars[pos-1]) != '\\' {
			value = append(value, chars[pos])
			token := Token(fmt.Sprintf("STRING %s %s", string(value), string(value[1:len(value)-1])))
			return pos + 1, append(tokens, token)
		}
	}
}
