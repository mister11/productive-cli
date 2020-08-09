package productive

import "time"

func formatDate(date time.Time) string {
	return date.Format("2006-01-02")
}