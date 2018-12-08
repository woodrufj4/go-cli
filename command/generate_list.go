package command

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"go-cli/structs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/cli"
)

type GenerateListCommand struct {
	Ui       cli.Ui
	FileName string
}

// Shows the commands help in the UI.
func (gl *GenerateListCommand) Help() string {
	helpText := `
Usage: go-cli generate list [options]

  This command generates a csv file based on the
  iTunes API top 100 list.
`
	return strings.TrimSpace(helpText)
}

// Shows the shortened help within a list of commands.
func (gl *GenerateListCommand) Synopsis() string {
	return "Generates csv based on top 100 albums"
}

// Run calls the iTunes api and generates a csv
// based on the data from the reponse.
func (gl *GenerateListCommand) Run(args []string) int {
	gl.Ui.Info("Running generate list command:")
	gl.Ui.Info("Calling iTunes API...")

	resp, err := http.Get("https://itunes.apple.com/us/rss/topalbums/limit=100/json")

	if gl.hasError(err, true) {
		return 1
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)

	if gl.hasError(readErr, true) {
		return 1
	}

	gl.Ui.Info("Parsing response...")
	var jsonResp structs.ApiResponse

	jsonErr := gl.parseJSON(body, &jsonResp)

	if gl.hasError(jsonErr, true) {
		gl.Ui.Error("Unabled to parse API response.")
		return 1
	}

	gl.Ui.Info("Generating csv...")
	genErr := gl.generateCSV(jsonResp)

	if gl.hasError(genErr, true) {
		return 1
	}

	gl.Ui.Info("Done.")

	return 0
}

// Checks if an error exists. You can optionally display the error
// in the UI.
func (gl *GenerateListCommand) hasError(err error, displayError bool) bool {
	if err == nil {
		return false
	}

	if displayError == true {
		gl.Ui.Error(err.Error())
	}

	return true
}

// Parses the data into JSON format.
func (gl *GenerateListCommand) parseJSON(data []byte, v interface{}) error {
	if json.Valid(data) != true {
		return errors.New("Invalid json format from API response.")
	}

	return json.Unmarshal(data, &v)
}

// Generates a csv file based on the JSON data.
func (gl *GenerateListCommand) generateCSV(jsonData structs.ApiResponse) error {

	filePtr, err := os.Create(time.Now().UTC().Format("20060102150405") + ".csv")
	// filePtr, err := os.Create(strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ".csv")
	defer filePtr.Close()

	if gl.hasError(err, false) {
		gl.Ui.Error("Could not create csv file:")
		return err
	}

	writer := csv.NewWriter(filePtr)
	defer writer.Flush()

	// File headers
	writeErr := writer.Write([]string{"iTunesID", "Category", "Name", "Artist", "Link", "Price", "ReleaseDate"})

	if gl.hasError(writeErr, false) {
		gl.Ui.Error("Could note write headers to csv file:")
		return writeErr
	}

	// File data
	for _, v := range jsonData.Feed.Entries {

		innerErr := writer.Write([]string{
			v.Id.Attributes.Id,
			v.Category.Attributes.Label,
			v.Name.Label,
			v.Artist.Label,
			v.Link.Attributes.Href,
			v.Price.Attributes.Amount,
			v.ReleaseDate.Label,
		})

		if gl.hasError(innerErr, false) {
			gl.Ui.Error("Could not complete writing dataset to file:")
			return innerErr
		}
	}
	return nil
}
