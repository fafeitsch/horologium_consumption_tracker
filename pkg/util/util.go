package util

import (
	"fmt"
	"time"
)

func FormatDate(year int, month int, day int) time.Time {
	result, _ := time.Parse(DateFormat, fmt.Sprintf("%04d-%02d-%02d", year, month, day))
	return result
}

func FormatDatePtr(year int, month int, day int) *time.Time {
	result := FormatDate(year, month, day)
	return &result
}

const DateFormat = "2006-01-02"
