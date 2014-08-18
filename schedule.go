package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"

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

// NewEventFromYAML creates an Event from a YAML file
func NewEventFromYAML(path string) (event *Event, err error) {
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

func (event *Event) processLinks() {
	ignorePrefix := "public\\/"
	for _, period := range event.Schedule {
		for _, item := range period.Items {
			filename := regexp.MustCompile(ignorePrefix + "(.+)\\.md")
			matches := filename.FindStringSubmatch(item.DeckMarkdownPath)
			fmt.Printf("%#v\n", filename)
			fmt.Printf("%#v\n", matches)
			if matches != nil {
				item.DeckHTMLPath = matches[1] + "/index.html"
			}
			fmt.Printf("%#v\n", item)
		}
	}
	// TODO: changes to items aren't available here
	fmt.Printf("%#v\n", event)
}

func (event *Event) generateHTML() (out string, err error) {
	html := `
<html>
  <head>
    <title>{{ .Title }}</title>
  </head>
  <body>
    <table id="main-details">
      <tr>
        <th>Event</th><td>Example 3 day training</td>
      </tr>
      <tr>
        <th>Location</th><td>Stark & Wayne HQ<br/>Buffalo, NY</td>
      </tr>
    </table>
    {{range .Schedule}}
      <h2>{{ .Label }}</h2>
      {{ if .Items }}
      <ul>
      {{ range .Items}}
				{{ if .Name }}
        <li class="item">
					{{ .Name }}
					{{ if .DeckMarkdownPath }}
						<a href="{{ .DeckHTMLPath }}">session slides</a>
					{{ end }}
					{{ if .LabMarkdownPath }}
						<a href="{{ .LabMarkdownPath }}">lab/workshop</a>
					{{ end }}
				</li>
				{{ else }}
				<li class="break" />
				{{ end }}
      {{ end }}
      </ul>
      {{ else }}
      <p>No items scheduled for today.</p>
      {{ end }}
    {{ else }}
    <p>No days scheduled yet.</p>
    {{end}}
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
	event, err := NewEventFromYAML(path)
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	event.processLinks()

	html, err := event.generateHTML()
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	println(html)
}
