package main

import (
	. "harvesterd"
	"flag"
	"fmt"
)

const version string = "0.0.1"

type Options struct {
	configFile string
	verbose    bool
	debug      bool
	help       bool
}

var options Options

func init() {
	flag.StringVar(&options.configFile, "config", "/etc/harvesterd.conf", "config filename")
	flag.BoolVar(&options.verbose, "verbose", false, "raise log level to verbose")
	flag.BoolVar(&options.debug, "debug", false, "raise log level to debug")
	flag.BoolVar(&options.help, "help", false, "help display this help")

	flag.Usage = help
}

func main() {
	flag.Parse()
	if options.help {
		help()
		return
	}

	run()
}

func help() {
	fmt.Printf("\033[1mharvesterd v%s\033[0m\n", version)
	fmt.Printf("Low footprint harvesterd and parser for events and logs\n")
	fmt.Printf("Máximo Cuadros Ortiz <mcuadros@gmail.com>\n\n")

	fmt.Printf("Usage:\n")
	flag.PrintDefaults()
}

func run() {
	harvesterd := NewHarvesterd()
	harvesterd.Configure(options.configFile)
	harvesterd.Boot()
	harvesterd.Run()
}
