package main

import (
	"fmt"
	"os"

	"go-cli/command"

	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	commandUI := &cli.BasicUi{
		Writer:      os.Stdout,
		Reader:      os.Stdin,
		ErrorWriter: os.Stderr,
	}

	commands := command.Commands(commandUI)

	cli := &cli.CLI{
		Name:     "go-cli",
		Version:  "0.1.1",
		Args:     args,
		Commands: commands,
	}

	exitCode, err := cli.Run()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Error executing CLI: %s\n", err.Error())
	}

	return exitCode
}
