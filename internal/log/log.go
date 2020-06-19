package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Debug(message interface{}) {
	color.White("%s", message)
}

func Info(format string, args ...interface{}) {
	color.Magenta(format, args)
}

func Error(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
}
