package command

import (
	"github.com/mitchellh/cli"
)

const (
	API_URL = "https://itunes.apple.com/us/rss/topalbums/limit=100/json"
)

func Commands(commandUI cli.Ui) map[string]cli.CommandFactory {

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
		"generate list-images": func() (cli.Command, error) {
			return &GenerateListImagesCommand{
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
