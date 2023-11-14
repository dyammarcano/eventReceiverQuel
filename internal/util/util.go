package util

import (
	"log"
	"os"
	"regexp"
)

// ExtractDigits extracts the digits from a string
func ExtractDigits(input string) string {
	re := regexp.MustCompile("[0-9]+")
	digits := re.FindAllString(input, -1)
	result := ""
	for _, digit := range digits {
		result += digit
	}
	return result
}

// CheckErr prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func CheckErr(msg any) {
	if msg != nil {
		log.Printf("Error: %v", msg)
		os.Exit(1)
	}
}

func CheckErrFatal(msg any) {
	if msg != nil {
		log.Fatalf("Error: %v", msg)
	}
}

// CheckErrPanic prints the msg with the prefix 'Error:' and panics. If the msg is nil, it does nothing.
func CheckErrPanic(msg any) {
	if msg != nil {
		log.Panicf("Error: %v", msg)
	}
}
