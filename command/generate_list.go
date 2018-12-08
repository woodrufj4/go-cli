package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type GenerateListCommand struct {
	Ui cli.Ui
}

func (gl *GenerateListCommand) Help() string {
	helpText := `
Usage: go-cli generate list [options]

  This command generates a csv file based on the
  iTunes API top 100 list.
`
	return strings.TrimSpace(helpText)
}

func (gl *GenerateListCommand) Synopsis() string {
	return "Generates csv based on top 100 albums"
}

func (gl *GenerateListCommand) Run(args []string) int {
	return 0
}
