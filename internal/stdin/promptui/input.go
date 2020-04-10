package promptui

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/mister11/productive-cli/internal/utils"
)

func (promptUiStdin *PromptUiStdin) InputMasked(label string) string {
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

func (promptUiStdin *PromptUiStdin) Input(label string) string {
	prompt := promptui.Prompt{
		Label: color.MagentaString(label),
	}

	result, err := prompt.Run()

	if err != nil {
		utils.ReportError("Error running prompt.", err)
	}
	return result
}

func (promptUiStdin *PromptUiStdin) InputMultiple(label string) []string {
	index := 1
	var inputs []string
	for isEnd := false; !isEnd; {
		input := promptUiStdin.Input(fmt.Sprintf("%s %d (empty to finish)", label, index))
		if len(input) == 0 {
			isEnd = true
			continue
		}
		inputs = append(inputs, input)
		index++
	}
	return inputs
}
