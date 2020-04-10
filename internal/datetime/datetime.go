package datetime

import (
	"github.com/jinzhu/now"
	"github.com/mister11/productive-cli/internal/utils"
	"time"
)

const yyyyMmDd = "2006-01-02"

func Now() time.Time {
	return time.Now()
}

func WeekStart() time.Time {
	return now.Monday()
}

func WeekEnd() time.Time {
	return now.Sunday().AddDate(0, 0, -2)
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
