package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo"
)

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
