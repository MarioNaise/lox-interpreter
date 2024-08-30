package lox

import "time"

func globals() map[string]callable {
	return map[string]callable{
		"clock": {evaluate: getTime, String: nativeFnString},
	}
}

func getTime() any           { return time.Now().Unix() }
func nativeFnString() string { return "<native fn>" }
