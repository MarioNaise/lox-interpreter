package lox

import (
	"time"
)

func globals() map[string]any {
	return map[string]any{
		"clock": &builtin{function: getTime},
	}
}

func getTime(*interpreter, []any) any { return float64(time.Now().Unix()) }
