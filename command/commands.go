package command

import (
	"os"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func Commands(commandUI cli.Ui) map[string]cli.CommandFactory {

	if terminal.IsTerminal(int(os.Stdout.Fd())) {

		commandUI = &cli.ColoredUi{
			OutputColor: cli.UiColorNone,
			InfoColor:   cli.UiColorBlue,
			ErrorColor:  cli.UiColorRed,
			WarnColor:   cli.UiColorYellow,
			Ui:          commandUI,
		}
	}

	all := map[string]cli.CommandFactory{
		"generate": func() (cli.Command, error) {
			return &GenerateCommand{
				Ui: commandUI,
			}, nil
		},
		"generate list": func() (cli.Command, error) {
			return &GenerateListCommand{
				Ui: commandUI,
			}, nil
		},
		"transfer": func() (cli.Command, error) {
			return &TransferCommand{
				Ui: commandUI,
			}, nil
		},
	}

	return all
}
