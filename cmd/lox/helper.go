package lox

import (
	"fmt"
	"os"
	"regexp"
)

func getFileContent(fileName string) string {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return string(fileContents)
}

func getPathFromFile(filePath string) string {
	r := regexp.MustCompile(`[^\/]+$`)
	dirPath := r.ReplaceAllString(filePath, "")
	return dirPath
}

func getFileName(filePath string) string {
	r := regexp.MustCompile(`[^\/]+$`)
	dirPath := r.FindString(filePath)
	return dirPath
}

func printErrors(errors []loxError) {
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}
}

func handleStmt(stmt stmt, i *interpreter) {
	defer continueOnError()
	i.resolveStmt(stmt)
	i.execute(stmt)
}

func handleExpr(exp expression, i *interpreter) {
	defer continueOnError()
	i.resolveExpr(exp)
	fmt.Println(i.stringify(i.evaluate(exp)))
}

func handleExprEval(exp expression, i *interpreter) {
	defer exitOnError()
	fmt.Println(i.stringify(i.evaluate(exp)))
}

func exitOnError() {
	if r := recover(); r != nil {
		fmt.Fprintln(os.Stderr, r)
		os.Exit(70)
	}
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

func recoverLoxError(t token) {
	if r := recover(); r != nil {
		switch r := r.(type) {
		case loxError:
			err := newError(r.message, t.line)
			panic(err)
		default:
			panic(r)
		}
	}
}
