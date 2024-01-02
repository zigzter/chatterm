package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}
