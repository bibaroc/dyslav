package main

import (
	"os"
	"strconv"
)

func envStr(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
func envBool(key string, def bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return def
		}
		return b
	}
	return def
}
