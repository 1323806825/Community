package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

	seed := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(seed)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[generator.Intn(len(charset))]
	}
	return string(b)
}
