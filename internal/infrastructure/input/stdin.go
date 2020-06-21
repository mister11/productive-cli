package input

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/mister11/productive-cli/internal/domain/config"
	"github.com/mister11/productive-cli/internal/utils"
)

type StdinPrompt struct{}

func NewStdinPrompt() *StdinPrompt {
	return &StdinPrompt{}
}


func (stdIn *StdinPrompt) Input(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: color.MagentaString(label),
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return result, nil
}

func (stdIn *StdinPrompt) InputMasked(label string) string {
	prompt := promptui.Prompt{
		Label: color.MagentaString(label),
		Mask:  '*',
	}

	result, err := prompt.Run()

	if err != nil {
		utils.ReportError("Error running prompt.", err)
	}
	return result
}

func (stdIn *StdinPrompt) InputMultiple(label string) []string {
	index := 1
	var inputs []string
	for isEnd := false; !isEnd; {
		input, _ := stdIn.Input(fmt.Sprintf("%s %d (empty to finish)", label, index))
		if len(input) == 0 {
			isEnd = true
			continue
		}
		inputs = append(inputs, input)
		index++
	}
	return inputs
}

func (stdIn *StdinPrompt) SelectOne(label string, options []interface{}) interface{} {
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
		utils.ReportError("Input failed", err)
	}

	return options[index]
}

func (stdIn *StdinPrompt) SelectOneWithSearch(label string, options []config.Project, searchFunction func(string, int) bool) interface{} {
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
