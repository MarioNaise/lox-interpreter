package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lox"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "tokenize":
		handleTokenizeCommand()
	case "parse":
		handleParseCommand()
	case "evaluate":
		handleEvaluateCommand()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func handleTokenizeCommand() {
	r := getFileReader(os.Args[2])
	ok := lox.Tokenize(r)
	if !ok {
		os.Exit(65)
	}
	os.Exit(0)
}

func handleParseCommand() {
	r := getFileReader(os.Args[2])
	ok := lox.Parse(r)
	if !ok {
		os.Exit(65)
	}
	os.Exit(0)
}

func handleEvaluateCommand() {
	r := getFileReader(os.Args[2])
	ok := lox.Evaluate(r)
	if !ok {
		os.Exit(65)
	}
	os.Exit(0)
}

func getFileReader(filename string) io.Reader {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return strings.NewReader(string(fileContents))
}
