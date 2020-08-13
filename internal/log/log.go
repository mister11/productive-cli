package log

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

func Debug(format string, args ...interface{}) {
	fmt.Println(aurora.White(fmt.Sprintf(format, args)))
}

func Info(format string, args ...interface{}) {
	fmt.Println(aurora.Cyan(fmt.Sprintf(format, args)))
}

func Error(format string, err ...error) {
	_ = fmt.Errorf(format, err)
}
