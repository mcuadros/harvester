package main

import collector "../src"

func main() {
	app := collector.Collector{}

	app.Configure()
	app.Boot()
	app.ReadFile()
}
