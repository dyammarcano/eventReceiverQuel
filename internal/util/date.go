package util

import (
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

func FormatDate(date time.Time) string {
	return date.Format(TimeFormat)
}

func DateISO8601(date time.Time) string {
	return date.Format(time.RFC3339)
}
