package main

import (
	. "collector"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "collector"
	app.Usage = "fight the loneliness!"
	app.Author = "Máximo Cuadros"
	app.Email = "mcuadros@gmail.com"

	app.Commands = []cli.Command{{
		Name:  "daemon",
		Usage: "add a task to the list",
		Flags: []cli.Flag{
			cli.BoolFlag{"verbose", "raise log level to info"},
			cli.BoolFlag{"debug", "raise log level to debug"},
			cli.StringFlag{"config, c", "/etc/collectord.conf", "config file"},
		},
		Action: daemon,
	}}

	app.Run(os.Args)
}

func daemon(c *cli.Context) {
	collector := NewCollector()
	collector.Configure(c.String("config"))
	collector.Boot()
	collector.Run()
}
