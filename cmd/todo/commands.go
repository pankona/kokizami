package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo"
)

// GlobalFlags can be used globally
var GlobalFlags = []cli.Flag{}

// Commands represents list of subcommands
var Commands = []cli.Command{
	{
		Name:   "start",
		Usage:  "start new task",
		Action: CmdStart,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "edit",
		Usage:  "edit task",
		Action: CmdEdit,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "list",
		Usage:  "show list of tasks",
		Action: CmdList,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "stop",
		Usage:  "stop task",
		Action: CmdStop,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound is called when specified subcommand is not found
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

// CmdStart starts a new task
func CmdStart(c *cli.Context) {
	for i, v := range c.Args() {
		switch i {
		case 0:
			t, err := todo.Start(v)
			if err != nil {
				log.Println(err)
			}
			log.Println(t)
		}
	}
}

// CmdEdit edits a specified task
func CmdEdit(c *cli.Context) {
	args := c.Args()
	if len(args) != 2 {
		log.Println("edit needs two arguments (id, desc)")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Println(err)
		return
	}
	desc := args[1]
	t, err := todo.Edit(id, desc)
	if err != nil {
		log.Println(err)
	}
	log.Println(t)
}

// CmdList shows ToDo list
func CmdList(c *cli.Context) {
	l, err := todo.List()
	if err != nil {
		log.Println(err)
		return
	}

	if len(l) == 0 {
		log.Println("list is empty")
		return
	}

	for _, v := range l {
		log.Println(v)
	}
}

// CmdStop update specified task's stopped_at
func CmdStop(c *cli.Context) {
	for i, v := range c.Args() {
		switch i {
		case 0:
			id, err := strconv.Atoi(v)
			if err != nil {
				log.Println(err)
				return
			}
			err = todo.Stop(id)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
