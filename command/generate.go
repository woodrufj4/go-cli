package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type GenerateCommand struct {
	Ui cli.Ui
}

func (g *GenerateCommand) Help() string {
	helpText := `
Usage: go-cli generate <subcommand> [options]

  This command groups commands for generating file(s).
  Users can generate a csv file of the top 100 albums or
  create a directory tree of album covers.

  Generate csv file of top 100

	$ go-cli generate list
	
  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(helpText)
}

func (g *GenerateCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (g *GenerateCommand) Synopsis() string {
	return "Generate file(s) based on iTunes API"
}
