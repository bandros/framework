package framework

import (
	"crypto/sha512"
	"encoding/hex"
)

func Password(plaintext string) string {
	h512 := sha512.New()
	h512.Write([]byte(plaintext))
	return hex.EncodeToString(h512.Sum(nil))
}

func ValidPassword(plaintext string, hash string) bool {
	if Password(plaintext) == hash {
		return true
	}

	return false
}
