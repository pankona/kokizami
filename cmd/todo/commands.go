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
	{
		Name:   "delete",
		Usage:  "delete task",
		Action: CmdDelete,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound is called when specified subcommand is not found
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

// CmdStart starts a new task
// todo start [new desc]
func CmdStart(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		log.Println("stop needs one arguments [id]")
		return
	}

	desc := args[0]
	t, err := todo.Start(desc)
	if err != nil {
		log.Println(err)
	}

	log.Println(t)
}

// CmdEdit edits a specified task
// todo edit [id] [new desc]
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
// todo list
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
// todo stop [id]
func CmdStop(c *cli.Context) {
	args := c.Args()
	switch len(args) {
	case 0:
		err := todo.StopAll()
		if err != nil {
			log.Println(err)
		}
	case 1:
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println(err)
			return
		}
		err = todo.Stop(id)
		if err != nil {
			log.Println(err)
		}
	default:
		log.Println("stop needs at most one arguments [id]")
		return
	}

}

// CmdDelete deletes specified task
// todo delete [id]
func CmdDelete(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		log.Println("delete needs one arguments [id]")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Println(err)
		return
	}

	err = todo.Delete(id)
	if err != nil {
		log.Println(err)
	}
}
