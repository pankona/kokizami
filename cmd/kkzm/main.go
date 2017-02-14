package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/pankona/kokizami"
)

func main() {
	// TODO: support multi platform
	kokizami.Initialize(os.Getenv("HOME") + "/.kokizami.db")

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "pankona"
	app.Email = "yosuke.akatsuka@gmail.com"
	app.Usage = "awesome task timer and tracker"
	app.Flags = GlobalFlags
	app.Action = CmdList // show list if no argument
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Run(os.Args)
}
