package datetime

import "time"

const yyyyMmDd = "2006-01-02"

type DateTimeProvider interface {
	Now() time.Time
	GetWeekDays() []time.Time
	ToISOTime(timeString string) time.Time
	Format(time time.Time) string
}
