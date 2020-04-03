package utils

import (
	"io/ioutil"
	"strings"
)

func ReadFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trimmedContent := strings.TrimSpace(string(content))
	return []byte(trimmedContent), nil
}
