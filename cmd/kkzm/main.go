package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pankona/kokizami"
	"github.com/urfave/cli"
)

func main() {
	// TODO: support multi platform
	err := kokizami.Initialize(os.Getenv("HOME") + "/.kokizami.db")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize: %v", err))
	}

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
	err = app.Run(os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
