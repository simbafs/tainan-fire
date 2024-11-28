package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func Digest(s string) string {
	r := sha256.Sum256([]byte(s))
	return hex.EncodeToString(r[:])
}
