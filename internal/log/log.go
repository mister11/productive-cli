package log

import (
	"github.com/fatih/color"
)

func Debug(message interface{}) {
	color.White("%s", message)
}

func Info(format string, args ...interface{}) {
	color.Magenta(format, args)
}

func Error(message string) {
	color.Red(message)
}
