package lox

type token struct {
	tokenType string
	lexeme    string
	literal   string
	line      int
}

func newToken(tokenType string, lexeme string, literal string, line int) token {
	return token{tokenType: tokenType, lexeme: lexeme, literal: literal, line: line}
}

func (t token) String() string {
	if t.literal == None {
		return t.tokenType + " " + t.lexeme
	} else {
		return t.tokenType + " " + t.lexeme + " " + t.literal
	}
}
