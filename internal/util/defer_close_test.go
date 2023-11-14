package util

import "testing"

func TestDeferClose(t *testing.T) {
	DeferClose(nil)
}
