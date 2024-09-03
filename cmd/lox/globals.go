package lox

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

func globals() map[string]any {
	return map[string]any{
		"read":   &builtin{function: readLn},
		"clock":  &builtin{function: getTime},
		"print":  &builtin{function: printLn, lenArgs: 1},
		"random": &builtin{function: random, lenArgs: 1},
		"sleep":  &builtin{function: sleep, lenArgs: 1},
		"import": &builtin{function: importFn, lenArgs: 1},
	}
}

func readLn(*interpreter, []any) any         { s := bufio.NewScanner(os.Stdin); s.Scan(); return s.Text() }
func getTime(*interpreter, []any) any        { return float64(time.Now().Unix()) }
func printLn(i *interpreter, args []any) any { fmt.Println(i.stringify(args[0])); return nil }

func random(_ *interpreter, args []any) any {
	if v, ok := args[0].(float64); ok && v > 0 {
		return float64(rand.Int64N(int64(v)))
	}
	// TODO: lox error
	return 0
}

func sleep(_ *interpreter, args []any) any {
	length, ok := args[0].(float64)
	if ok {
		time.Sleep(time.Duration(length) * time.Millisecond)
	}
	return nil
}

func importFn(i *interpreter, args []any) any {
	if filePath, ok := args[0].(string); ok {
		prevIndex := i.index
		defer func() { i.index = prevIndex }()
		i.index = i.index + getPathFromFile(filePath)
		content := getFileContent(joinBaseAndFilePath(i.index, filePath))
		p := newParser(content)
		p.parse()
		if len(p.parseErrors) == 0 {
			i.interpret(p.program)
		} else {
			fmt.Fprintln(os.Stderr, "--------")
			fmt.Fprintf(os.Stderr, "Error in %s:\n", filePath)
			printErrors(p.parseErrors)
			fmt.Fprintln(os.Stderr, "--------")

		}
	}
	return nil
}

func getFileContent(fileName string) string {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		// TODO: lox error
		return ""
	}
	return string(fileContents)
}

func joinBaseAndFilePath(base string, filePath string) string {
	if filePath[0] == '/' {
		return filePath
	}
	return base + getFileName(filePath)
}
