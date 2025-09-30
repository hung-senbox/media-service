package helper

import (
	"fmt"
	"time"
)

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func ValidateDateRange(start, end time.Time) bool {
	return !start.After(end)
}

func FormatRemainingDays(days int) string {
	if days <= 0 {
		return "Expired"
	}
	return fmt.Sprintf("%d", days)
}
