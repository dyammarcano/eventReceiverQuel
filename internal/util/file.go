package util

import (
	"os"
	"strings"
)

func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func RemoveExtension(filepath string) string {
	return filepath[:len(filepath)-len(GetExtension(filepath))]
}

func GetExtension(filepath string) string {
	if strings.Contains(filepath, ".") == true {
		return filepath[len(filepath)-4:]
	}
	return ""
}
