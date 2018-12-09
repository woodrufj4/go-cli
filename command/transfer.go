package command

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
	"github.com/mitchellh/cli"
)

type TransferCommand struct {
	Ui          cli.Ui
	Server      string
	FilePath    string
	NewFileName string
	Username    string
	Password    string
}

func (t *TransferCommand) Help() string {
	helpText := `
Usage: go-cli transfer <file-path> [options]

  This command transfers a file to an ftp server.

  Transfer a file:

	$ go-cli transfer ./some/file.csv

	General Options:

	-server=<string>
	  Sets the ftp server url.
	  Default: test.rebex.net:21

	-user=<string>
	  Sets username used to authenticate with the ftp server
	  
	-password=<string>
	  Sets the password used to authenticate with the ftp server

	-filename=<string>
	  Sets the file name to save as on the ftp server.
	  Default: The file name of the existing file.
`
	return strings.TrimSpace(helpText)
}

func (t *TransferCommand) Synopsis() string {
	return "Transfers a file to ftp server."
}

// Runs the transfer command to transfer a file
// to a remote ftp server.
func (t *TransferCommand) Run(args []string) int {

	if len(args) == 0 {
		return cli.RunResultHelp
	}

	if err := t.parseArgs(args); err != nil {
		t.Ui.Error("Could not parse arguments:")
		t.Ui.Output(err.Error())
		return 1
	}

	if err := t.verifyCredentials(t.Username, t.Password); err != nil {
		t.Ui.Error(err.Error())
		return 1
	}

	t.Ui.Info("Attempting to transfer file...")

	if err := t.transfer(); err != nil {
		t.Ui.Error(err.Error())
		return 1
	}

	t.Ui.Info("Transfer complete.")

	return 0
}

// Parses the command line arguments
func (t *TransferCommand) parseArgs(args []string) error {

	t.FilePath = path.Clean(args[0])

	if _, err := os.Stat(t.FilePath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("The file '%s' does not exist.", t.FilePath))
	}

	fs := flag.NewFlagSet("transfer general options", flag.ContinueOnError)
	fs.StringVar(&t.Server, "server", "test.rebex.net:21", "Sets ftp server url.")
	fs.StringVar(&t.Username, "user", "", "Sets the authN username.")
	fs.StringVar(&t.Password, "password", "", "Sets the authN password.")
	fs.StringVar(&t.NewFileName, "filename", filepath.Base(t.FilePath), "Sets the new filename to save as on the ftp server.")
	return fs.Parse(args[1:])
}

// Verifies that UserInfo credentials are set.
func (t *TransferCommand) verifyCredentials(identity string, secret string) error {
	if identity == "" || secret == "" {
		return errors.New("Both credentials are needed to authenticate with the ftp server.")
	}

	return nil
}

// Attemps to transfer the file to the ftp server
func (t *TransferCommand) transfer() error {

	conn, connErr := ftp.Dial(t.Server)

	if connErr != nil {
		t.Ui.Error(fmt.Sprintf("Could not connect to the ftp server '%s':", t.Server))
		return connErr
	}

	defer conn.Quit()

	if loginErr := conn.Login(t.Username, t.Password); loginErr != nil {
		t.Ui.Error(fmt.Sprintf("Could not login to the ftp server '%s':", t.Server))
		return loginErr
	}

	filePtr, fileErr := os.Open(t.FilePath)

	if fileErr != nil {
		t.Ui.Error(fmt.Sprintf("Could not open the file '%s':", t.FilePath))
		return fileErr
	}

	defer filePtr.Close()

	buff := bufio.NewReader(filePtr)

	if storErr := conn.Stor(t.NewFileName, buff); storErr != nil {
		t.Ui.Error(fmt.Sprintf("Could not transfer file to the ftp server '%s':", t.Server))
		return storErr
	}
	return nil
}
