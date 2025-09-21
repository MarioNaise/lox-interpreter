package main

import (
	"fmt"
	"os"

	"lox/cmd/lox"
)

func main() {
	if len(os.Args) == 1 {
		lox.Repl()
		return
	}

	if len(os.Args) == 2 {
		if ok := lox.Run(os.Args[1]); !ok {
			os.Exit(65)
		}
		return
	}

	command := os.Args[1]
	fileName := os.Args[2]

	handlers := map[string]func(string) bool{
		"tokenize": lox.Tokenize,
		"parse":    lox.Parse,
		"lint":     lox.Lint,
		"evaluate": lox.Evaluate,
		"run":      lox.Run,
	}

	if handler, ok := handlers[command]; ok {
		ok := handler(fileName)
		if !ok {
			os.Exit(65)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
