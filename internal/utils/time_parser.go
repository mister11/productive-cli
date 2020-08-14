package utils

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var TimeRegex = regexp.MustCompile(`^(?:(\d+)[:])?(\d+)$`)

// parses time of format HH:mm to minutes
// needed for the Productive API which accepts minutes
func ParseTime(time string) (string, error) {
	matches := TimeRegex.FindStringSubmatch(time)
	if len(matches) != 3 {
		return "", errors.New("wrong time format - only minutes or HH:mm format allowed")
	}
	hoursString := matches[1]
	minutesString := matches[2]
	if len(hoursString) == 0 {
		return strconv.Itoa(getMinutes(minutesString)), nil
	}
	return strconv.Itoa(getHours(hoursString)*60 + getMinutes(minutesString)), nil
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
