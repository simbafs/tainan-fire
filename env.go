package main

import (
	"os"
	"strconv"
)

func Getenv(key string, init string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return init
}

func GetenvInt64(key string, init int64) int64 {
	e := Getenv(key, "")
	if e == "" {
		return init
	}

	i, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}
