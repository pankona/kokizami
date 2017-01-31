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

	for _, v := range l {
		if v.IsDone() {
			continue
		}
		log.Println(v)
	}
}
