package utils

import (
	"github.com/mister11/productive-cli/internal/infrastructure/log"
)

func ReportError(message string, err error) {
	log.Error(message, err)
	panic(err)
}
