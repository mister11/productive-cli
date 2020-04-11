package utils

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"syscall"
	"testing"
)

func TestReadFileSpaces(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	content := []byte("   Test string    ")
	WriteFile(f.Name(), content)

	actualData, err := ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []byte("Test string"), actualData)
}

func TestReadFileNoSpaces(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	expectedData := []byte("Test string")
	WriteFile(f.Name(), expectedData)

	actualData, err := ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	assert.Equal(t, expectedData, actualData)
}

func TestReadFileMissing(t *testing.T) {
	result, err := ReadFile("non_existing_path")
	defer syscall.Unlink("non_existing_path")
	assert.Nil(t, result)
	assert.Error(t, err)
}