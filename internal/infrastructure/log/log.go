package log

import (
	"fmt"
	"github.com/fatih/color"
)

func Debug(format string, args ...interface{}) {
	color.White(format, args)
}

func Info(format string, args ...interface{}) {
	color.Magenta(format, args)
}

func Error(format string, err ...error) {
	_ = fmt.Errorf(format, err)
}
