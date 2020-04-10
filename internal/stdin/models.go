package stdin

import "github.com/mister11/productive-cli/internal/config"

type Stdin interface {
	InputMasked(label string) string
	Input(label string) string
	InputMultiple(label string) []string
	SelectOne(label string, options []interface{}) interface{}
	SelectOneWithSearch(label string, options []config.Project, searchFunction func(string, int) bool) interface{}
}
