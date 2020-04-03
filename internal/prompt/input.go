package prompt

import (
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"gitlab.com/mister11/productive-cli/internal/utils"
)

func InputMasked(message string) string {
	prompt := promptui.Prompt{
		Label: color.MagentaString(message),
		Mask:  '*',
	}

	result, err := prompt.Run()

	if err != nil {
		utils.ReportError("Error running prompt.", err)
	}
	return result
}

func Input(message string) string {
	prompt := promptui.Prompt{
		Label: color.MagentaString(message),
	}

	result, err := prompt.Run()

	if err != nil {
		utils.ReportError("Error running prompt.", err)
	}
	return result
}

func Confirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()

	return err == nil
}
