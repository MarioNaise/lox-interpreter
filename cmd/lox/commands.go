package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	PROMPT = "> "
	EXIT   = ".exit"
)

func Repl() {
	s := bufio.NewScanner(os.Stdin)
	i := newInterpreter(nil)
	fmt.Print(PROMPT)
	for s.Scan() {
		if s.Text() == EXIT {
			return
		}
		r := strings.NewReader(s.Text())
		i.parser = newParser(r)
		stmts, parseErrors := i.parse()
		if len(parseErrors) == 0 {
			for _, stmt := range stmts {
				switch stmt := stmt.(type) {
				case *stmtExpr:
					i.execute(&stmtPrint{stmt})
				default:
					handleStmt(stmt, i)
				}
			}
		}
		for _, err := range parseErrors {
			fmt.Fprintln(os.Stderr, replError(err.String()))
		}
		fmt.Print(PROMPT)
	}
}

func handleStmt(stmt stmtInterface, i *interpreter) {
	defer continueOnError()
	i.execute(stmt)
}

func continueOnError() {
	if r := recover(); r != nil {
		switch r := r.(type) {
		case loxError:
			fmt.Fprintln(os.Stderr, replError(r.String()))
		default:
			fmt.Fprintln(os.Stderr, r)
		}
	}
}

func replError(err string) string {
	reg := regexp.MustCompile(`\[line \d+\]\s`)
	return reg.ReplaceAllString(err, "")
}

func Tokenize(r io.Reader) bool {
	s := newScanner(r)
	tokens, errs := s.tokenize()
	for _, token := range tokens {
		fmt.Println(token)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Parse(r io.Reader) bool {
	p := newParser(r)
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

func Evaluate(r io.Reader) bool {
	defer exitOnError()
	i := newInterpreter(r)
	i.tokenize()
	expr := i.expression()
	if expr == nil {
		return true
	}
	errs := append(i.scanErrors, i.parseErrors...)
	if len(errs) == 0 {
		stmt := &stmtPrint{&stmtExpr{initializer: expr}}
		i.visitPrintStmt(stmt)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Run(r io.Reader) bool {
	i := newInterpreter(r)
	defer exitOnError()
	stmts, errs := i.parse()
	if len(errs) == 0 {
		i.interpret(stmts)
		return true
	}
	printErrors(errs)
	return false
}

func exitOnError() {
	if r := recover(); r != nil {
		fmt.Fprintln(os.Stderr, r)
		os.Exit(70)
	}
}

func printErrors(errors []loxError) {
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}
}
