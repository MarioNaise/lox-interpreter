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
	i := newInterpreter("")
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
					handleStmt(&stmtPrint{stmt.initializer}, i)
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

func Tokenize(str string) bool {
	s := newScanner(str)
	tokens, errs := s.tokenize()
	for _, token := range tokens {
		fmt.Println(token)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Parse(str string) bool {
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

func Evaluate(str string) bool {
	defer exitOnError()
	i := newInterpreter(str)
	i.tokenize()
	expr := i.expression()
	errs := append(i.scanErrors, i.parseErrors...)
	if expr == nil {
		printErrors(errs)
		return len(errs) == 0
	}
	if len(errs) == 0 {
		stmt := &stmtPrint{expr}
		i.visitPrintStmt(stmt)
	}
	printErrors(errs)
	return len(errs) == 0
}

func Run(str string) bool {
	i := newInterpreter(str)
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
