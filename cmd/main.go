package main

import (
	"fmt"
	"lox/cmd/lox"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		lox.Repl()
		return
	}
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
	case "run":
		handleRunCommand()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func handleTokenizeCommand() {
	r := getFileContent(os.Args[2])
	ok := lox.Tokenize(r)
	if !ok {
		os.Exit(65)
	}
}

func handleParseCommand() {
	r := getFileContent(os.Args[2])
	ok := lox.Parse(r)
	if !ok {
		os.Exit(65)
	}
}

func handleEvaluateCommand() {
	r := getFileContent(os.Args[2])
	ok := lox.Evaluate(r)
	if !ok {
		os.Exit(65)
	}
}

func handleRunCommand() {
	r := getFileContent(os.Args[2])
	ok := lox.Run(r)
	if !ok {
		os.Exit(65)
	}
}

func getFileContent(filename string) string {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return string(fileContents)
}
