package datetime

import (
	"github.com/jinzhu/now"
	"gitlab.com/mister11/productive-cli/internal/utils"
	"time"
)

const yyyyMmDd = "2006-01-02"

func Now() time.Time {
	return time.Now()
}

func NowFormatted() string {
	return time.Now().Format(yyyyMmDd)
}

func WeekStart() time.Time {
	return now.Monday()
}

func WeekEnd() time.Time {
	return now.Sunday().AddDate(0, 0, -2)
}

func MonthStartFormatted() string {
	return now.BeginningOfMonth().Format(yyyyMmDd)
}

func MonthEndFormatted() string {
	return now.EndOfMonth().Format(yyyyMmDd)
}

func YearStartFormatted() string {
	return now.BeginningOfYear().Format(yyyyMmDd)
}

func YearEndFormatted() string {
	return now.EndOfYear().Format(yyyyMmDd)
}

func ToISODate(dateString string) time.Time {
	time, err := time.Parse(yyyyMmDd, dateString)
	if err != nil {
		utils.ReportError("Error parsing "+dateString, err)
	}
	return time
}

func Format(time time.Time) string {
	return time.Format(yyyyMmDd)
}
