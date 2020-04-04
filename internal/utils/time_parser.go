package utils

import (
	"errors"
	"github.com/mister11/productive-cli/internal/log"
	"regexp"
	"strconv"
)

var timeRegex = regexp.MustCompile("^(?:(\\d+)[:])?(\\d+)$"				)

// parses time of format HH:mm to minutes
// needed for API which accepts minutes
func ParseTime(time string) int {
	matches := timeRegex.FindStringSubmatch(time)
	if len(matches) == 3 {
		return getHours(matches[1])*60 + getMinutes(matches[2])
	}
	if len(matches) == 2 {
		return getMinutes(matches[1])
	}
	log.Error("Wrong time format. You can enter either only minutes or HH:mm format")
	panic(errors.New("wrong time format"))
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