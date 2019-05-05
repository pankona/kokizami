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

	kkzm := &kokizami.Kokizami{}

	app.Before = func(ctx *cli.Context) error {
		u, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %v", err)
		}

		configDir := filepath.Join(u.HomeDir, ".config", "kokizami")
		err = os.MkdirAll(configDir, 0755) // #nosec
		if err != nil {
			return fmt.Errorf("failed to create directory on %v", configDir)
		}

		kkzm.DBPath = filepath.Join(configDir, "db")

		err = kkzm.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize kokizami: %v", err)
		}

		app.Metadata["kkzm"] = kkzm
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		return kkzm.Finalize()
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
