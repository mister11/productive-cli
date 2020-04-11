package utils

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"syscall"
	"testing"
)

func TestWriteFile(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())
	expectedData := []byte("Test string")
	WriteFile(f.Name(), expectedData)

	actualData, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, expectedData, actualData)
}

func TestWriteFileMissing(t *testing.T) {
	content := []byte("Test string")
	assert.Panics(t, func() { WriteFile("", content) })
}
