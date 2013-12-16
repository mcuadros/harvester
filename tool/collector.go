package main

import collector "github.com/mcuadros/collector/src"

func main() {
	app := collector.Collector{}

	app.Configure()
	app.Boot()
	app.ReadFile()
}
