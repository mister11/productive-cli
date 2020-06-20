package utils

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/mister11/productive-cli/internal/log"
)

var timeRegex = regexp.MustCompile(`^(?:(\d+)[:])?(\d+)$`)

// parses time of format HH:mm to minutes
// needed for the Productive API which accepts minutes
func ParseTime(time string) string {
	matches := timeRegex.FindStringSubmatch(time)
	if len(matches) != 3 {
		log.Error("Wrong time format. You can enter either only minutes or HH:mm format", nil)
		panic(errors.New("wrong time format"))
	}
	hoursString := matches[1]
	minutesString := matches[2]
	if len(hoursString) == 0 {
		return strconv.Itoa(getMinutes(minutesString))
	}
	return strconv.Itoa(getHours(hoursString)*60 + getMinutes(minutesString))
}

func getHours(hoursString string) int {
	hours, err := strconv.Atoi(hoursString)
	if err != nil {
		panic(err)
	}
	return hours
}

func getMinutes(minutesString string) int {
	minutes, err := strconv.Atoi(minutesString)
	if err != nil {
		panic(err)
	}
	return minutes
}

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}