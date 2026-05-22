package vm

type tokenType int

type token struct {
	lexeme string
	tokenType
	line int
}

const (
	// Single-character tokens.
	tokenLeftParen tokenType = iota
	tokenRightParen
	tokenLeftBrace
	tokenRightBrace
	tokenComma
	tokenDot
	tokenMinus
	tokenPlus
	tokenSemicolon
	tokenSlash
	tokenStar
	// One or two character tokens.
	tokenBang
	tokenBangEqual
	tokenEqual
	tokenEqualEqual
	tokenGreater
	tokenGreaterEqual
	tokenLess
	tokenLessEqual
	// Literals.
	tokenIdentifier
	tokenString
	tokenNumber
	// Keywords.
	tokenAnd
	tokenClass
	tokenElse
	tokenFalse
	tokenFor
	tokenFun
	tokenIf
	tokennil
	tokenOr
	tokenPrint
	tokenReturn
	tokenSuper
	tokenThis
	tokenTrue
	tokenVar
	tokenWhile
	tokenError
	tokenEOF
)
