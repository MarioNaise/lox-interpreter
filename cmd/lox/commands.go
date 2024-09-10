package lox

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const (
	PROMPT = "> "
	EXIT   = ".exit"
)

func Repl() {
	s := bufio.NewScanner(os.Stdin)
	i := newInterpreter("", "")
	for fmt.Print(PROMPT); s.Scan(); fmt.Print(PROMPT) {
		if s.Text() == EXIT {
			return
		}
		i.parser = newParser(s.Text())
		stmts, parseErrors := i.parse()
		if len(parseErrors) == 0 {
			for _, stmt := range stmts {
				switch stmt := stmt.(type) {
				case *stmtExpr:
					handleExpr(stmt.initializer, i)
				default:
					handleStmt(stmt, i)
				}
			}
		}
		for _, err := range parseErrors {
			fmt.Fprintln(os.Stderr, replError(err.String()))
		}
	}
}

func replError(err string) string {
	reg := regexp.MustCompile(`\[line \d+\]\s`)
	return reg.ReplaceAllString(err, "")
}

func Tokenize(filePath string) bool {
	str := getFileContent(filePath)
	s := newScanner(str)
	tokens, errs := s.tokenize()
	for _, token := range tokens {
		fmt.Println(token)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Parse(filePath string) bool {
	str := getFileContent(filePath)
	p := newParser(str)
	stmts, errs := p.parse()
	aP := astPrinter{}
	if len(stmts) == 1 &&
		len(errs) == 1 &&
		errs[0].message ==
			"Expected ';' after expression." {
		aP.print(stmts)

		return true
	}
	if len(errs) == 0 {
		aP.print(stmts)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Evaluate(filePath string) bool {
	defer exitOnError()
	str := getFileContent(filePath)
	dirPath := getPathFromFile(filePath)
	i := newInterpreter(str, dirPath)
	i.tokenize()
	expr := i.expression()
	errs := append(i.scanErrors, i.parseErrors...)
	if expr == nil {
		printErrors(errs)
		return len(errs) == 0
	}
	if len(errs) == 0 {
		handleExprEval(expr, i)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Run(filePath string) bool {
	defer exitOnError()
	str := getFileContent(filePath)
	dirPath := getPathFromFile(filePath)
	i := newInterpreter(str, dirPath)
	stmts, errs := i.parse()
	i.resolver.resolve(stmts)
	if len(errs) == 0 {
		i.interpret(stmts)
		return true
	}
	printErrors(errs)
	return false
}
