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
	ok := lox.Tokenize(os.Args[2])
	if !ok {
		os.Exit(65)
	}
}

func handleParseCommand() {
	ok := lox.Parse(os.Args[2])
	if !ok {
		os.Exit(65)
	}
}

func handleEvaluateCommand() {
	ok := lox.Evaluate(os.Args[2])
	if !ok {
		os.Exit(65)
	}
}

func handleRunCommand() {
	ok := lox.Run(os.Args[2])
	if !ok {
		os.Exit(65)
	}
}
