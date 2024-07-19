package utils

import (
	"time"
)

func FormatDate(date time.Time) string {
	return date.Format(time.RFC3339)
}
