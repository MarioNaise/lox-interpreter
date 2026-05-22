package vm

type scanner struct {
	start   int
	current int
	line    int
	src     string
}

func newScanner(src string) *scanner {
	return &scanner{0, 0, 1, src}
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *scanner) scanToken() token {
	s.skipWhitespace()
	s.start = s.current

	if s.isAtEnd() {
		return token{"", tokenEOF, s.line}
	}

	c := s.advance()
	switch {
	case c == '(':
		return s.makeToken(tokenLeftParen)
	case c == ')':
		return s.makeToken(tokenRightParen)
	case c == '{':
		return s.makeToken(tokenLeftBrace)
	case c == '}':
		return s.makeToken(tokenRightBrace)
	case c == ';':
		return s.makeToken(tokenSemicolon)
	case c == ',':
		return s.makeToken(tokenComma)
	case c == '.':
		return s.makeToken(tokenDot)
	case c == '-':
		return s.makeToken(tokenMinus)
	case c == '+':
		return s.makeToken(tokenPlus)
	case c == '/':
		return s.makeToken(tokenSlash)
	case c == '*':
		return s.makeToken(tokenStar)
	case c == '!':
		if s.match('=') {
			return s.makeToken(tokenBangEqual)
		} else {
			return s.makeToken(tokenBang)
		}
	case c == '=':
		if s.match('=') {
			return s.makeToken(tokenBangEqual)
		} else {
			return s.makeToken(tokenBang)
		}
	case c == '<':
		if s.match('=') {
			return s.makeToken(tokenBangEqual)
		} else {
			return s.makeToken(tokenBang)
		}
	case c == '>':
		if s.match('=') {
			return s.makeToken(tokenBangEqual)
		} else {
			return s.makeToken(tokenBang)
		}
	case c == '"':
		return s.string()
	case isDigit(c):
		return s.number()
	case isAlpha(c):
		return s.identifier()
	default:
		return s.errorToken("Unexpected character.")
	}
}

func (s *scanner) makeToken(t tokenType) token {
	return token{s.src[s.start:s.current], t, s.line}
}

func (s *scanner) errorToken(msg string) token {
	return token{msg, tokenError, s.line}
}

func (s *scanner) skipWhitespace() {
	for !s.isAtEnd() {
		c := s.peek()
		switch c {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *scanner) identifierType() tokenType {
	switch s.src[s.start] {
	case 'a':
		return s.checkKeyword(1, 2, "nd", tokenAnd)
	case 'c':
		return s.checkKeyword(1, 4, "lass", tokenClass)
	case 'e':
		return s.checkKeyword(1, 3, "lse", tokenElse)
	case 'f':
		if s.current-s.start > 1 {
			switch s.src[s.start+1] {
			case 'a':
				return s.checkKeyword(2, 3, "lse", tokenFalse)
			case 'o':
				return s.checkKeyword(2, 1, "r", tokenFor)
			case 'u':
				return s.checkKeyword(2, 1, "n", tokenFun)
			}
		}
		return tokenIdentifier
	case 'i':
		return s.checkKeyword(1, 1, "f", tokenIf)
	case 'n':
		return s.checkKeyword(1, 2, "il", tokennil)
	case 'o':
		return s.checkKeyword(1, 1, "r", tokenOr)
	case 'p':
		return s.checkKeyword(1, 4, "rint", tokenPrint)
	case 'r':
		return s.checkKeyword(1, 5, "eturn", tokenReturn)
	case 's':
		return s.checkKeyword(1, 4, "uper", tokenSuper)
	case 't':
		if s.current-s.start > 1 {
			switch s.src[s.start+1] {
			case 'h':
				return s.checkKeyword(2, 2, "is", tokenThis)
			case 'r':
				return s.checkKeyword(2, 2, "ue", tokenTrue)
			}
		}
		return tokenIdentifier
	case 'v':
		return s.checkKeyword(1, 2, "ar", tokenVar)
	case 'w':
		return s.checkKeyword(1, 4, "hile", tokenWhile)
	default:
		return tokenIdentifier
	}
}

func (s *scanner) checkKeyword(start int, length int, rest string, t tokenType) tokenType {
	if s.current-s.start == start+length && s.src[s.start+start:s.start+start+length] == rest {
		return t
	}
	return tokenIdentifier
}

func (s *scanner) identifier() token {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}
	return s.makeToken(s.identifierType())
}

func (s *scanner) number() token {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the ".".
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	return s.makeToken(tokenNumber)
}

func (s *scanner) string() token {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}

	// The closing quote.
	s.advance()
	return s.makeToken(tokenString)
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *scanner) advance() byte {
	s.current++
	return s.src[s.current-1]
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.src[s.current]
}

func (s *scanner) peekNext() byte {
	if s.current+1 >= len(s.src) {
		return 0
	}
	return s.src[s.current+1]
}

func (s *scanner) match(expected byte) bool {
	if s.isAtEnd() || s.src[s.current] != expected {
		return false
	}
	s.current++
	return true
}
