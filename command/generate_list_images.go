package command

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go-cli/structs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/cli"
)

type GenerateListImagesCommand struct {
	Ui       cli.Ui
	FileName string
	Path     string
}

// Shows the commands help in the UI.
func (gli *GenerateListImagesCommand) Help() string {
	helpText := `
Usage: go-cli generate list-images [options]

  This command generates a csv file of album image links 
  based on the iTunes API top 100 list.

  General Options:

  -filename=<string>
	Sets the file name for the generated csv. The extension
	will be added at runtime.
	Example: kung-fu
	
  -path=<string>
	Sets the path for where the generated csv will be placed.
	Example: ~/some/path/
`
	return strings.TrimSpace(helpText)
}

// Shows the shortened help within a list of commands.
func (gli *GenerateListImagesCommand) Synopsis() string {
	return "Generates csv of album image links based on top 100 albums"
}

// Run calls the iTunes api and generates a csv
// of image links based on the data from the reponse.
func (gli *GenerateListImagesCommand) Run(args []string) int {

	if err := gli.parseFlags(args); err != nil {
		gli.Ui.Error("Unable to parse flag options:")
		gli.Ui.Output(err.Error())
		return 1
	}

	if err := gli.verifyPath(gli.Path); err != nil {
		gli.Ui.Error("File path error:")
		gli.Ui.Output(err.Error())
		return 1
	}

	gli.Ui.Info("Running generate list-images command...")
	gli.Ui.Info("Calling iTunes API...")

	resp, err := http.Get(API_URL)

	if err != nil {
		gli.Ui.Error("API call error:")
		gli.Ui.Output(err.Error())
		return 1
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		gli.Ui.Error("API reponse error:")
		gli.Ui.Output(err.Error())
		return 1
	}

	var jsonResp structs.ApiResponse

	if err := gli.parseJSON(body, &jsonResp); err != nil {
		gli.Ui.Error("JSON parse error:")
		gli.Ui.Output(err.Error())
		return 1
	}

	gli.Ui.Info("Generating csv...")
	if err := gli.generateCSV(jsonResp); err != nil {
		gli.Ui.Error("Could not generate csv file:")
		gli.Ui.Output(err.Error())
		return 1
	}

	gli.Ui.Info("Done.")
	return 0
}

// Parse the flag options
func (gli *GenerateListImagesCommand) parseFlags(args []string) error {
	fs := flag.NewFlagSet("list general options", flag.ContinueOnError)
	fs.StringVar(&gli.FileName, "filename", "", "Sets the filename for the generated csv.")
	fs.StringVar(&gli.Path, "path", "", "Sets the path for where the generated csv will be placed.")
	return fs.Parse(args)
}

// Verifies that the path is a directory
func (gli *GenerateListImagesCommand) verifyPath(path string) error {
	if path == "" {
		return nil
	}

	stat, statErr := os.Stat(path)

	if statErr != nil {
		return statErr
	}

	if stat.IsDir() == false {
		return errors.New("The provided path is not a directory")
	}

	return nil
}

// Parses the data into JSON format.
func (gli *GenerateListImagesCommand) parseJSON(data []byte, v interface{}) error {
	if json.Valid(data) != true {
		return errors.New("Invalid json format from API response.")
	}

	return json.Unmarshal(data, &v)
}

// Generates the csv file
func (gli *GenerateListImagesCommand) generateCSV(jsonData structs.ApiResponse) error {

	if gli.FileName == "" {
		gli.FileName = fmt.Sprintf("%s_list_images", time.Now().UTC().Format("20060102150405"))
	}

	if gli.Path != "" {
		gli.Path = gli.Path + "/"
	}

	filePath := path.Clean(fmt.Sprintf("%s%s.csv", gli.Path, gli.FileName))

	filePtr, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer filePtr.Close()
	gli.Ui.Warn(fmt.Sprintf("Generated file: %s", filePath))

	writer := csv.NewWriter(filePtr)
	defer writer.Flush()

	// File headers
	if err := writer.Write([]string{"iTunesID", "smallImageLocation", "mediumImageLocation", "largeImageLocation"}); err != nil {
		return err
	}

	// File data
	for _, v := range jsonData.Feed.Entries {

		var smallImageLocation, mediumImageLocation, largeImageLocation string

		// Taking the naive approach here and assuming the
		// order of the small, medium, and large image locations
		// will always be the same.
		for k, i := range v.Images {
			switch k {
			case 0:
				smallImageLocation = i.Label

			case 1:
				mediumImageLocation = i.Label

			case 2:
				largeImageLocation = i.Label
			}
		}

		err := writer.Write([]string{
			v.Id.Attributes.Id,
			smallImageLocation,
			mediumImageLocation,
			largeImageLocation,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
