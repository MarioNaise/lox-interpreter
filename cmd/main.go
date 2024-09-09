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
	if len(os.Args) == 2 {
		handleRunCommand(os.Args[1])
		return
	}
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	fileName := os.Args[2]

	switch command {
	case "tokenize":
		handleTokenizeCommand(fileName)
	case "parse":
		handleParseCommand(fileName)
	case "evaluate":
		handleEvaluateCommand(fileName)
	case "run":
		handleRunCommand(fileName)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func handleTokenizeCommand(fileName string) {
	ok := lox.Tokenize(fileName)
	if !ok {
		os.Exit(65)
	}
}

func handleParseCommand(fileName string) {
	ok := lox.Parse(fileName)
	if !ok {
		os.Exit(65)
	}
}

func handleEvaluateCommand(fileName string) {
	ok := lox.Evaluate(fileName)
	if !ok {
		os.Exit(65)
	}
}

func handleRunCommand(fileName string) {
	ok := lox.Run(fileName)
	if !ok {
		os.Exit(65)
	}
}
