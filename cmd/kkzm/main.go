package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pankona/kokizami"
	"github.com/urfave/cli"
)

func main() {
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

	app.Before = func(ctx *cli.Context) error {
		u, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %v", err)
		}

		kkzm := &kokizami.Kokizami{}
		err = kkzm.Initialize(filepath.Join(u.HomeDir, ".config", "kokizami", "db"))
		if err != nil {
			return fmt.Errorf("failed to initialize: %v", err)
		}

		app.Metadata["kkzm"] = kkzm
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
