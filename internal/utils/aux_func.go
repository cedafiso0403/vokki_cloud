package utils

import (
	"database/sql"
	"time"
)

func FormatDate(date time.Time) string {
	return date.Format(time.RFC3339)
}

func ConvertNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
