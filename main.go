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
	Title    string `yaml:"title"`
	Location string `yaml:"location"`
	Schedule []*struct {
		Label string `yaml:"label"`
		Items []*struct {
			Name             string `yaml:"name"`
			DeckMarkdownPath string `yaml:"deck"`
			DeckHTMLPath     string
			LabMarkdownPath  string `yaml:"lab"`
			LabHTMLPath      string
		}
	}
	Students []*struct {
		Name    string
		Email   string
		Login   string `yaml:"login"`
		Host    string `yaml:"host"`
		SshPort int    `yaml:"port"`
	}
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
			if matches != nil {
				item.DeckHTMLPath = "/" + matches[1] + "/index.html"
			}
			matches = filename.FindStringSubmatch(item.LabMarkdownPath)
			if matches != nil {
				item.LabHTMLPath = "/labs#!" + matches[1] + ".md"
			}
		}
	}
}

func (event *Event) processStudentsForSharedServer(host string, port int) {
	for index, student := range event.Students {
		student.Login = fmt.Sprintf("student%d", index+1)
		student.Host = host
		student.SshPort = port
	}
}

func (event *Event) generateHTML() (out string, err error) {
	html := `
  <!DOCTYPE html>
<html lang="en">
  <head>
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- Bootstrap -->
    <link href="bootstrap-3.2.0-dist/css/bootstrap.min.css" rel="stylesheet">
  </head>
  <body>
    <div class="container-fluid">
      <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8">
          <h1>{{ .Title }}</h1>
          <table class="table">
            <tr>
              <th>Event</th><td>AT&amp;T Testing Training</td>
            </tr>
            <tr>
              <th>Location</th><td>{{ .Location }}</td>
            </tr>
          </table>
        </div>
        <div class="col-md-2"></div>
      </div>
      <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8">
          <h2>Schedule</h2>
          <table class="table table-hover">
            <tr><th>Day</th><th>Topics</th></tr>
            {{ range .Schedule }}
            <tr>
              <td>{{ .Label }}</td>
              <td>
                <ul>
                {{ range .Items }}
                  {{ if .Name }}
                  <li class="item">
                    {{ .Name }}
                    {{ if .DeckMarkdownPath }}
                      [<a href="{{ .DeckHTMLPath }}">session slides</a>]
                    {{ end }}
                    {{ if .LabMarkdownPath }}
                      [<a href="{{ .LabHTMLPath }}">lab/workshop</a>]
                    {{ end }}
                  </li>
                  {{ else }}
                    <li class="break" />
                  {{ end }}
                {{ end }}
                </ul>
              </td>
            </tr>
            {{ end }}

            </table>
          </div>
          <div class="col-md-2"></div>

        </div>
      </div>
      <br/>
      <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8">
          <h2>Students</h2>

          <table class="table table-hover">
            <tr><th>Name</th><th>Email</th><th>Login</th></tr>

            {{ range .Students }}
            <tr>
            <td>{{ .Name }}</td>
            <td>{{ .Email }}</td>
            <td>{{ .Login }}</td>
            </tr>
            {{ end }}

            </table>
          </div>
          <div class="col-md-2"></div>
        </div>
      </div>

      <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8">
          <h3>Email all students</h3>
          <p>
            Select All & Copy into clipboard.
          </p>
          <input type="text" value="{{ range .Students }}{{ if .Email }}{{ .Name }} <{{ .Email }}>, {{ end }}{{ end }}" size=120>
        </div>
        <div class="col-md-2"></div>
      </div>
    </div>
    <br/>
    <br/>
    <br/>

        <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="bootstrap-3.2.0-dist/js/bootstrap.min.js"></script>
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
	var host string
	var port int
	flag.StringVar(&host, "host", "", "host for shared student server")
	flag.IntVar(&port, "port", 22, "ssh port for shared student server")
	flag.Parse()

	path := flag.Arg(0)
	event, err := NewEventFromYAML(path)
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	event.processLinks()
	event.processStudentsForSharedServer(host, port)

	html, err := event.generateHTML()
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	println(html)
}
