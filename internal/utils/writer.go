package utils

import (
	"io/ioutil"
)

func WriteFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0644)
}
