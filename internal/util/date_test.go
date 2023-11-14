package util

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	date := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	if FormatDate(date) != "2019-01-01 00:00:00" {
		t.Errorf("FormatDate(date) is not 2019-01-01 00:00:00")
		return
	}
}
