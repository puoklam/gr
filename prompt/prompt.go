package prompt

import "github.com/manifoldco/promptui"

func Run(label any, validate promptui.ValidateFunc) (string, error) {
	pt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	return pt.Run()
}
