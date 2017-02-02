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
		Name:   "add",
		Usage:  "add new task",
		Action: CmdAdd,
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
		Name:   "done",
		Usage:  "mark specified task as done",
		Action: CmdDone,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound is called when specified subcommand is not found
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

// CmdAdd adds a new todo
func CmdAdd(c *cli.Context) {
	for i, v := range c.Args() {
		switch i {
		case 0:
			t, err := todo.Add(v)
			if err != nil {
				log.Println(err)
			}
			log.Println(t)
		}
	}
}

// CmdEdit adds a new todo
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

	for _, v := range l {
		if v.IsDone() {
			continue
		}
		log.Println(v)
	}
}

// CmdDone marks specified ToDo as done
func CmdDone(c *cli.Context) {
	for i, v := range c.Args() {
		switch i {
		case 0:
			id, err := strconv.Atoi(v)
			if err != nil {
				log.Println(err)
				return
			}
			err = todo.Done(id)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
