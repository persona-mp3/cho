package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func generateHash(name string) string {
	shaBytes :=  sha256.Sum256([]byte(name))
	return hex.EncodeToString(shaBytes[:])
}

