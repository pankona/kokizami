package main

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo"
)

// CmdList shows ToDo list
func CmdList(c *cli.Context) {
	l, err := todo.List()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(l)
}
