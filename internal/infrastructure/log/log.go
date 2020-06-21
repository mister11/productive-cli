package log

import (
	"fmt"
	"github.com/fatih/color"
)

func Debug(message interface{}) {
	color.White("%s", message)
}

func Info(format string, args ...interface{}) {
	color.Magenta(format, args)
}

func Error(format string, err ...error) {
	_ = fmt.Errorf(format, err)
}
