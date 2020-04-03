package prompt

import (
	"github.com/manifoldco/promptui"
	"gitlab.com/mister11/productive-cli/internal/config"
	"gitlab.com/mister11/productive-cli/internal/utils"
)

func SelectOne(label string, options []interface{}) interface{} {
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
		utils.ReportError("Prompt failed", err)
	}

	return options[index]
}

func SelectOneWithSearch(label string, options []config.Project, searchFunction func(string, int) bool) interface{} {
	if len(options) == 0 {
		return nil
	}
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
