package log

import (
	"github.com/fatih/color"
)

func Debug(message interface{}) {
	color.White("%s", message)
}

func Info(message string) {
	color.Magenta(message)
}

func Error(message string) {
	color.Red(message)
}
