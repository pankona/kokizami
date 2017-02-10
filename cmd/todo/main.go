package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo"
)

func main() {
	todo.Initialize(nil, os.Getenv("HOME")+"/.todo.db")

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "pankona"
	app.Email = "yosuke.akatsuka@gmail.com"
	app.Usage = "awesome todo manager"
	app.Flags = GlobalFlags
	app.Action = CmdList // show list if no argument
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Run(os.Args)
}
