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
	}
}

func readLn(*interpreter, []any) any         { s := bufio.NewScanner(os.Stdin); s.Scan(); return s.Text() }
func getTime(*interpreter, []any) any        { return float64(time.Now().Unix()) }
func printLn(i *interpreter, args []any) any { fmt.Println(i.stringify(args[0])); return nil }

func random(_ *interpreter, args []any) any {
	if v, ok := args[0].(float64); ok && v > 0 {
		return float64(rand.Int64N(int64(v)))
	}
	return 0
}

func sleep(_ *interpreter, args []any) any {
	length, ok := args[0].(float64)
	if ok {
		time.Sleep(time.Duration(length) * time.Millisecond)
	}
	return nil
}
