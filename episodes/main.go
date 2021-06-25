package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type Episode struct {
	Number       int
	Date         time.Time
	Eu           bool
	Title        string
	YouTube      string
	GuestName    string
	GuestURL     string
	HostName     string
	HostURL      string
	ShowNotesURL string
}

func main() {

	var spreadsheetFilename string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "filename",
				Aliases:     []string{"f"},
				Value:       "planning.csv",
				Usage:       "Episode planning spreadsheet in CSV format",
				Destination: &spreadsheetFilename,
			},
		},
		Action: func(c *cli.Context) error {
			return episode(spreadsheetFilename)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func episode(spreadsheetFilename string) error {
	readme := "README.md"
	t, err := getTemplate(readme, "000-template")
	if err != nil {
		return fmt.Errorf("getting template from 000-template/README.md")
	}

	spreadsheet, err := os.Open(spreadsheetFilename)
	if err != nil {
		return fmt.Errorf("reading spreadsheet file: %v", err)
	}

	ep, err := readEpisode(spreadsheet)
	if err != nil {
		return fmt.Errorf("episode info error: %v", err)
	}

	// Create directory for this episode number
	dir := fmt.Sprintf("%03d", ep.Number)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModeDir|0755)
		if err != nil {
			return fmt.Errorf("making %s directory: %v", dir, err)
		}
	}

	filename := fmt.Sprintf("%s/%s", dir, readme)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("writing to file %s: %v", filename, err)
	}
	defer f.Close()

	// Write the template with the episode data
	err = t.Execute(f, ep)
	if err != nil {
		return fmt.Errorf("executing template: %v", err)
	}

	return nil
}

func getTemplate(readme string, templatePath string) (*template.Template, error) {
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("reading ./%s directory, are you in the right place?", templatePath)
	}

	templateText, err := ioutil.ReadFile(templatePath + "/" + readme)
	if err != nil {
		return nil, fmt.Errorf("reading template file: %v", err)
	}

	funcMap := template.FuncMap{
		"upcoming": func(i int) int {
			return i + 2
		},
		"upcomingDate": func(d time.Time) time.Time {
			return d.AddDate(0, 0, 14)
		},
		"friendly": func(d time.Time) string {
			return d.Format("2 January 2006")
		},
		"iso": func(d time.Time) string {
			return d.Format("20060102")
		},
	}
	return template.Must(template.New("README").Funcs(funcMap).Parse(string(templateText))), nil
}

func readEpisode(in io.Reader) (*Episode, error) {
	var p promptui.Prompt

	p = promptui.Prompt{Label: "Spreadsheet row number"}
	row, _ := p.Run()
	rowNum, err := strconv.Atoi(row)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(in)
	fields := []string{}
	for rowNum > 0 {
		rowNum--
		fields, err = r.Read()
		if err != nil {
			return nil, fmt.Errorf("reading spreadsheet row: %v", err)
		}
	}

	ep := Episode{
		Title: fields[3],
	}

	ep.Date, err = time.Parse("2-January-2006", fields[0])
	if err != nil {
		return nil, fmt.Errorf("failed parsing date from %s: %v", fields[0], err)
	}

	ep.Number, err = strconv.Atoi(fields[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse episode number from %s: %v", fields[1], err)
	}

	host := ""
	hostURL := ""
	switch fields[4] {
	case "Liz":
		host = "Liz Rice"
		hostURL = "https://twitter.com/lizrice"
	case "Duffie":
		host = "Duffie Cooley"
		hostURL = "https://twitter.com/mauilion"
	default:
		p = promptui.Prompt{
			Label: "Host name",
		}

		host, err = p.Run()
		if err != nil {
			return nil, err
		}

		p = promptui.Prompt{
			Label: "Host URL",
		}

		hostURL, err = p.Run()
		if err != nil {
			return nil, err
		}
	}

	ep.HostName = host
	ep.HostURL = hostURL

	p = promptui.Prompt{
		Label:   "Guest",
		Default: fields[5],
	}

	ep.GuestName, err = p.Run()
	if err != nil {
		return nil, err
	}

	p = promptui.Prompt{
		Label:   "Guest URL",
		Default: "https://",
	}

	ep.GuestURL, err = p.Run()
	if err != nil {
		return nil, err
	}

	p = promptui.Prompt{
		Label:   "YouTube URL",
		Default: "https://",
	}

	ep.YouTube, err = p.Run()
	if err != nil {
		return nil, err
	}

	ep.ShowNotesURL = fmt.Sprintf("/episodes/%03d", ep.Number)
	switch fields[2] {
	case "EU":
		ep.Eu = true
	case "US":
		ep.Eu = false
	default:
		return nil, fmt.Errorf("is this EU or US: %s", fields[2])
	}

	return &ep, nil
}
