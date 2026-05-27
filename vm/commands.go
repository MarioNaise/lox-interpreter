// Package vm implements the virtual machine for the Lox programming language.
package vm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Run(args []string) {
	switch len(args) {
	case 0:
		repl()
	case 1:
		runFile(args[0])
	default:
		fmt.Fprintln(os.Stderr, "Usage: lox vm [file]")
		os.Exit(64)
	}
}

const (
	replPrompt = "> "
	replExit   = ".exit"
)

func repl() {
	vm := new(vm)
	s := bufio.NewScanner(os.Stdin)
	for fmt.Print(replPrompt); s.Scan(); fmt.Print(replPrompt) {
		text := s.Text()
		if strings.HasPrefix(text, replExit) {
			return
		}
		vm.interpret(text)
	}
}

func runFile(filePath string) {
	vm := new(vm)
	src, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	result := vm.interpret(string(src))
	if result == interpretCompileError {
		os.Exit(65)
	}
	if result == interpretRuntimeError {
		os.Exit(70)
	}
}
