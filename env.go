package main

import "os"

func Getenv(key string, init string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return init
}
