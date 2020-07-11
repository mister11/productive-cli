package input

import "github.com/mister11/productive-cli/internal/domain"

type Prompt interface {
	Input(label string) (string, error)
	InputMasked(label string) (string, error)
	InputMultiline(label string) ([]string, error)
	SelectOne(label string, options []interface{}) (interface{}, error)
	SelectOneWithSearch(label string, options []domain.TrackedProject, searchFunction func(string, int) bool) interface{}
}