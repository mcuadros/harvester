package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"

	. "github.com/mcuadros/harvester/src"
)

// Define a template.
const configurationTemplate = `[{{.Name}}{{if .AllowMultiple}} "name"{{end}}]{{range .Fields}}
-{{.Name}} (type:{{.Type}}{{if .Default}}, default: {{.Default}}{{end}}): {{.Description}}{{end}}

`

type Options struct {
	configFile string
	verbose    bool
	debug      bool
	help       bool
	spec       string
}

var version string
var options Options

func init() {
	flag.StringVar(&options.configFile, "config", "/etc/harvester.conf", "config filename")
	flag.BoolVar(&options.verbose, "verbose", false, "raise log level to verbose")
	flag.BoolVar(&options.debug, "debug", false, "raise log level to debug")
	flag.BoolVar(&options.help, "help", false, "display this help")
	flag.StringVar(&options.spec, "spec", "", "display the specs for any config group, 'all' returns all the groups")

	flag.Usage = help
}

func main() {
	flag.Parse()

	switch {
	case options.help:
		help()
	case len(options.spec) != 0:
		spec(options.spec)
	default:
		run()
	}
}

func help() {
	fmt.Printf("\033[1mharvester (%s)\033[0m\n", version)
	fmt.Printf("Low footprint collector and parser for events and logs\n")
	fmt.Printf("Máximo Cuadros Ortiz <mcuadros@gmail.com>\n\n")

	fmt.Printf("Usage:\n")
	flag.PrintDefaults()
}

func spec(group string) {
	template := template.Must(template.New("tmpl").Parse(configurationTemplate))

	any := false
	if group == "all" {
		any = true
	}

	definitions := GetConfig().GetDescription()
	found := false
	for _, definition := range definitions {
		if definition.Name == group || any {
			err := template.Execute(os.Stdout, definition)
			if err != nil {
				fmt.Fprintf(os.Stderr, "executing template:", err)
			}

			found = true
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "unable to find spec for '%s'", group)
	}
}

func run() {
	harvester := NewHarvester()
	harvester.Configure(options.configFile)
	harvester.Boot()
	harvester.Run()
}
