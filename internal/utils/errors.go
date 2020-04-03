package utils

import (
	"gitlab.com/mister11/productive-cli/internal/log"
)

func ReportError(message string, err error) {
	log.Error(message)
	panic(err)
}
