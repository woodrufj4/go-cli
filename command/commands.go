package command

import (
	"github.com/mitchellh/cli"
)

func Commands(commandUI cli.Ui) map[string]cli.CommandFactory {

	coloredUI := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorBlue,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,
		Ui:          commandUI,
	}

	all := map[string]cli.CommandFactory{
		"generate": func() (cli.Command, error) {
			return &GenerateCommand{
				Ui: coloredUI,
			}, nil
		},
		"generate list": func() (cli.Command, error) {
			return &GenerateListCommand{
				Ui: coloredUI,
			}, nil
		},
		"transfer": func() (cli.Command, error) {
			return &TransferCommand{
				Ui: coloredUI,
			}, nil
		},
	}

	return all
}
