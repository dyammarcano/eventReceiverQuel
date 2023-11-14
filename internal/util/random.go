package util

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// GetRandomId returns a random uuid as a string
func GetRandomId() string {
	return uuid.New().String()
}

// VerifyId verifies that the given string is a valid uuid
func VerifyId(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// GenerateRandomString generates a random string of the given length
func GenerateRandomString(length, split int) string {
	characters := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	if split == -1 {
		randomChars := make([]byte, length)
		for i := 0; i < length; i++ {
			randomIndex := random.Intn(len(characters))
			randomChars[i] = characters[randomIndex]
		}
		return string(randomChars)
	}

	randomChars := make([]byte, length+(length-1)/split)
	for i, j := 0, 0; i < length; i, j = i+1, j+1 {
		if i > 0 && i%split == 0 {
			randomChars[j] = '-'
			j++
		}
		randomIndex := random.Intn(len(characters))
		randomChars[j] = characters[randomIndex]
	}
	return string(randomChars)
}

// RandomNumber generates a random number of the given length
func RandomNumber(length int) string {
	characters := "0123456789"
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomChars := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := random.Intn(len(characters))
		randomChars[i] = characters[randomIndex]
	}
	return string(randomChars)
}
