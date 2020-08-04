package input

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/mister11/productive-cli/internal/domain"
)

type StdinPrompt struct{}

func NewStdinPrompt() StdinPrompt {
	return StdinPrompt{}
}

func (stdIn StdinPrompt) Input(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: color.MagentaString(label),
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return result, nil
}

func (stdIn StdinPrompt) InputMasked(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: color.MagentaString(label),
		Mask:  '*',
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return result, nil
}

func (stdIn StdinPrompt) InputMultiline(label string) ([]string, error) {
	index := 1
	var inputs []string
	for isEnd := false; !isEnd; {
		input, err := stdIn.Input(fmt.Sprintf("%s %d (empty to finish)", label, index))
		if err != nil {
			return nil, err
		}
		if len(input) == 0 {
			isEnd = true
			continue
		}
		inputs = append(inputs, input)
		index++
	}
	return inputs, nil
}

func (stdIn StdinPrompt) SelectOne(label string, options []interface{}) (interface{}, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
		Templates: &promptui.SelectTemplates{
			Active:   "\U0001F872 {{ .Name | cyan }}",
			Inactive: "{{ .Name }}",
			Selected: "\U0001F872 {{ .Name | cyan }}",
		},
	}

	index, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return options[index], nil
}

func (stdIn StdinPrompt) SelectOneWithSearch(
	label string,
	options []domain.TrackedProject,
	searchFunction func(string, int) bool,
) interface{} {
	prompt := promptui.Select{
		Label: label,
		Items: options,
		Templates: &promptui.SelectTemplates{
			Active:   "\U0001F872 {{ .DealName | cyan }} - {{ .ServiceName | cyan}}",
			Inactive: "{{ .DealName }} - {{ .ServiceName }}",
			Selected: "\U0001F872 {{ .DealName | cyan }} - {{ .ServiceName | cyan}}",
		},
		Searcher: searchFunction,
	}

	index, _, err := prompt.Run()

	if err != nil {
		return nil
	}

	return options[index]
}
