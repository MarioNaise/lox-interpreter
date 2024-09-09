package lox

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

func globals() map[string]any {
	return map[string]any{
		"read":     &builtin{function: readLn},
		"clock":    &builtin{function: getTime},
		"print":    &builtin{function: printLn, lenArgs: 1},
		"random":   &builtin{function: random, lenArgs: 1},
		"sleep":    &builtin{function: sleep, lenArgs: 1},
		"string":   &builtin{function: stringify, lenArgs: 1},
		"parseNum": &builtin{function: parseNum, lenArgs: 1},
		"load":     &builtin{function: load, lenArgs: 1},
	}
}

func readLn(*interpreter, []any, token) any {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return s.Text()
}

func getTime(*interpreter, []any, token) any          { return float64(time.Now().Unix()) }
func printLn(i *interpreter, args []any, t token) any { fmt.Println(i.stringify(args[0])); return nil }

func random(_ *interpreter, args []any, t token) any {
	if v, ok := args[0].(float64); ok && v > 0 {
		return float64(rand.Int64N(int64(v)))
	}
	err := newError("random - Argument must be a positive number.", t.line)
	panic(err)
}

func sleep(_ *interpreter, args []any, t token) any {
	length, ok := args[0].(float64)
	if ok {
		time.Sleep(time.Duration(length) * time.Millisecond)
	} else {
		err := newError("sleep - Argument must be a number.", t.line)
		panic(err)
	}
	return nil
}

func stringify(i *interpreter, args []any, t token) any {
	return i.stringify(args[0])
}

func parseNum(i *interpreter, args []any, t token) any {
	str, ok := args[0].(string)
	if !ok {
		err := newError("parseNum - Argument must be a string.", t.line)
		panic(err)
	}
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		err := newError("parseNum - Number couldn't be parsed.", t.line)
		panic(err)
	}
	return num
}

func load(i *interpreter, args []any, t token) any {
	if filePath, ok := args[0].(string); ok {
		prevIndex := i.index
		defer func() {
			i.index = prevIndex
			if r := recover(); r != nil {
				switch r := r.(type) {
				case returnValue:
					return
				default:
					// TODO: clean up
					fmt.Fprintln(os.Stderr, "--------")
					fmt.Fprintf(os.Stderr, "Error in %s:\n", filePath)
					fmt.Fprintln(os.Stderr, r)
					fmt.Fprintln(os.Stderr, "--------")
					panic(r)
				}
			}
		}()
		i.index = i.index + getPathFromFile(filePath)
		content := getFileContentLoad(joinBaseAndFilePath(i.index, filePath), t)
		p := newParser(content)
		p.parse()
		if len(p.parseErrors) == 0 {
			// TODO: resolve before interpreting in calling file
			i.resolver.resolve(p.program)
			i.interpret(p.program)
		} else {
			// TODO: clean up
			fmt.Fprintln(os.Stderr, "--------")
			fmt.Fprintf(os.Stderr, "Error in %s:\n", filePath)
			printErrors(p.parseErrors)
			fmt.Fprintln(os.Stderr, "--------")

		}
	}
	return nil
}

func getFileContentLoad(fileName string, t token) string {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		err := newError("Could not read file "+fileName, t.line)
		panic(err)
	}
	return string(fileContents)
}

func joinBaseAndFilePath(base string, filePath string) string {
	if filePath[0] == '/' {
		return filePath
	}
	return base + getFileName(filePath)
}
