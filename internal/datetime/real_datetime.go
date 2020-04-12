package datetime

import (
	"github.com/jinzhu/now"
	"github.com/mister11/productive-cli/internal/utils"
	"time"
)

type RealDateTimeProvider struct {}

func NewRealTimeDateProvider() *RealDateTimeProvider {
	return &RealDateTimeProvider{}
}

func (dateTimeProvider *RealDateTimeProvider) Now() time.Time {
	return time.Now()
}

func (dateTimeProvider *RealDateTimeProvider) WeekStart() time.Time {
	return now.Monday()
}

func (dateTimeProvider *RealDateTimeProvider) WeekEnd() time.Time {
	return now.Sunday().AddDate(0, 0, -2)
}

func (dateTimeProvider *RealDateTimeProvider) ToISOTime(timeString string) time.Time {
	time, err := time.Parse(yyyyMmDd, timeString)
	if err != nil {
		utils.ReportError("Error parsing "+timeString, err)
	}
	return time
}

func (dateTimeProvider *RealDateTimeProvider) Format(time time.Time) string {
	return time.Format(yyyyMmDd)
}
