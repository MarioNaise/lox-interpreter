package vm

import "fmt"

func compile(src string) {
	s := newScanner(src)
	line := -1
	for t := s.scanToken(); t.tokenType != tokenEOF; t = s.scanToken() {
		if t.line != line {
			fmt.Printf("%4d ", t.line)
			line = t.line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%s'\n", t.tokenType, t.lexeme)
	}
}
