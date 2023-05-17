package prompt

import (
	"io"

	"github.com/manifoldco/promptui"
)

var p = promptui.Prompt{
	Label:    "",
	Validate: nil,
}

func SetStdin(rc io.ReadCloser) {
	p.Stdin = rc
}

func Run(label any, validate promptui.ValidateFunc) (string, error) {
	p.Label = label
	p.Validate = validate
	return p.Run()
}
