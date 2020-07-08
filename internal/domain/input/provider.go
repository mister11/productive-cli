package input

import (
	"github.com/mister11/productive-cli/internal/domain/config"
)

type Prompt interface {
	Input(label string) (string, error)
	InputMasked(label string) (string, error)
	InputMultiline(label string) ([]string, error)
	SelectOne(label string, options []interface{}) interface{}
	SelectOneWithSearch(label string, options []config.Project, searchFunction func(string, int) bool) interface{}
}