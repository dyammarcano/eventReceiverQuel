package util

import (
	"encoding/hex"
	"strconv"
	"strings"
)

// Byte2HexStr converts a byte array to a hex string
func Byte2HexStr(ba []byte) string {
	var sb strings.Builder
	for _, b := range ba {
		sb.WriteString(strconv.FormatInt(int64(b), 16))
	}
	return sb.String()
}

// HexStr2Byte converts a hex string to a byte array
func HexStr2Byte(ba string) []byte {
	binary := make([]byte, len(ba)/2)
	i := 0
	for i < len(ba) {
		num, _ := strconv.ParseInt(ba[i:i+2], 16, 64)
		binary[i/2] = byte(num)
		i += 2
	}
	return binary
}

func StringToHex(input string) string {
	return strings.ToUpper(hex.EncodeToString([]byte(input)))
}

func HexToString(input string) string {
	decoded, _ := hex.DecodeString(input)
	return string(decoded)
}
