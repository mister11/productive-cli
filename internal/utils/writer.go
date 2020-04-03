package utils

import (
	"io/ioutil"
)

func WriteFile(path string, content []byte) {
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		ReportError("Error writing to file "+path, err)
	}
}
