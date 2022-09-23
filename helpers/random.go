package helpers

import (
	"math/rand"
	"time"
)

// RandAlphanumericString ...
func RandAlphanumericString(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return StringWithCharset(length, charset)
}

// RandLowerAlphanumericString ...
func RandLowerAlphanumericString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	return StringWithCharset(length, charset)
}

// StringWithCharset ...
func StringWithCharset(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}