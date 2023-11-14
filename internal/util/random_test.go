package util

import "testing"

func TestGetRandomId(t *testing.T) {
	id := GetRandomId()
	if !VerifyId(id) {
		t.Errorf("GetRandomId() = %s; want a valid uuid", id)
	}
}

func TestGenerateRandomString(t *testing.T) {
	serialKey := GenerateRandomString(35, 7)

	if len(serialKey) != 39 {
		t.Errorf("GenerateRandomString() = %s; want a string of length 35", serialKey)
	}
}
