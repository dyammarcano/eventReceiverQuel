package util

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

// ReturnFixString returns a string with a fixed size and adds "..." if the string is bigger than the size
func ReturnFixString(str string, size int) string {
	if len(str) > size {
		return str[0:size] + "..."
	} else {
		return str
	}
}

// Upper returns the string in uppercase
func Upper(data string) string {
	return strings.ToUpper(data)
}

// SanitizeString removes any special characters that could cause issues
func SanitizeString(str string) string {
	// Remove any special characters that could cause issues
	regex := regexp.MustCompile(`[^\w\-.,:;()<>{}\[\]@!#$%^&*_+=|\\/~?]`)
	sanitizedStr := regex.ReplaceAllString(str, "")

	// Remove any leading or trailing whitespace
	sanitizedStr = strings.TrimSpace(sanitizedStr)

	return sanitizedStr
}

// Sha256Hash returns a sha256 hash of a string
func Sha256Hash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}

// Tokenizer returns a token of a string with a given length
func Tokenizer(str string, length int) string {
	return Sha256Hash(str)[0:length]
}
