package main

import (
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v1"
)

// Event contains an entire training schedule
type Event struct {
	Title    string        `yaml:"title"`
	Location string        `yaml:"location"`
	Schedule []EventPeriod `yaml:"schedule"`
}

// EventPeriod contains the schedule items for a particular time period, such as a day
type EventPeriod struct {
	Label string         `yaml:"label"`
	Items []ScheduleItem `yaml:"items"`
}

// ScheduleItem describes a scheduled item
type ScheduleItem struct {
	Name             string `yaml:"name"`
	DeckMarkdownPath string `yaml:"deck"`
	DeckHTMLPath     string
	LabMarkdownPath  string `yaml:"lab"`
	LabHTMLPath      string
}

func importEvent(path string) (event Event, err error) {
	file, err := os.Open(path)
	if err != nil {
		// TODO: how to wrap error with context?
		println("File does not exist:", err.Error())
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		// TODO: how to wrap error with context?
		println("Could not read file: ", err.Error())
		return
	}

	err = yaml.Unmarshal(data, &event)
	if err != nil {
		// TODO: how to wrap error with context?
		println("Could not unmarshall YAML: ", err.Error())
		return
	}
	return
}

func generateHTML(event Event) (out string, err error) {
	html := `
<html>
  <head>
    <title>{{ .Title }}</title>
  </head>
  <body>

  </body>
</html>
`
	tmpl, err := template.New("schedule").Parse(html)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, event)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func main() {
	flag.Parse()

	// fmt.Sprintf("%#v\n", flag.Args())
	// fmt.Sprintf("%#v\n", flag.Arg(0))
	// fmt.Sprintf("%#v\n", flag.Arg(1))

	path := flag.Arg(0)
	event, err := importEvent(path)
	if err != nil {
		println("Error: " + err.Error())
		return
	}
	html, err := generateHTML(event)
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	println(html)
}
