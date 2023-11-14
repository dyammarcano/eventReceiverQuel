package util

import (
	"fmt"
	"testing"
)

func TestHex(t *testing.T) {
	message := []byte("hello world")

	bytes := Byte2HexStr(message)

	if bytes != "68656c6c6f20776f726c64" {
		t.Errorf("Expected 68656c6c6f20776f726c64, got %s", bytes)
	}
}

func TestHexStr2Byte(t *testing.T) {
	message := "68656c6c6f20776f726c64"

	bytes := HexStr2Byte(message)

	if string(bytes) != "hello world" {
		t.Errorf("Expected hello world, got %s", string(bytes))
	}
}

func TestStringToHex(t *testing.T) {
	input := "ARQUIBM;;,05000.00-7"
	hexString := StringToHex(input)
	fmt.Printf("length: %d, string: %s\n", len(hexString), hexString)
}

func TestHexToString(t *testing.T) {
	input := "52494F424D2C30353030302E30302D37"
	hexToString := HexToString(input)
	fmt.Printf("length: %d, hexToString: %s\n", len(hexToString), hexToString)
}
