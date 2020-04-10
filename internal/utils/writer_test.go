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
	content := []byte("Test string")
	WriteFile(f.Name(), content)

	writtenData, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, content, writtenData)
}

func TestWriteFileMissing(t *testing.T) {
	content := []byte("Test string")
	WriteFile("nonexistingfile", content)
	defer syscall.Unlink("nonexistingfile")
	assert.Panics(t, func() { WriteFile("", content) })
}
