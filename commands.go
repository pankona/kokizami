package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "add",
		Usage:  "add new task",
		Action: command.CmdAdd,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
