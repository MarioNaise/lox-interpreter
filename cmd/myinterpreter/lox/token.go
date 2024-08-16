package lox

type token struct {
	tokenType string
	lexeme    string
	literal   string
}

func newToken(tokenType string, lexeme string, literal string) token {
	return token{tokenType: tokenType, lexeme: lexeme, literal: literal}
}

func (t token) String() string {
	return t.tokenType + " " + t.lexeme + " " + t.literal
}
