package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo"
)

func main() {
	todo.Initialize(nil)

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "pankona"
	app.Email = "yosuke.akatsuka@gmail.com"
	app.Usage = "awesome todo manager"
	app.Flags = GlobalFlags
	app.Action = func(c *cli.Context) error {
		log.Println("action here!")
		return cli.NewExitError("action failed", 3)
	}
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Run(os.Args)
}
