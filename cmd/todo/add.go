package main

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo"
)

// CmdAdd adds a new todo
func CmdAdd(c *cli.Context) {
	for i, v := range c.Args() {
		switch i {
		case 0:
			t := todo.NewToDo(v)
			err := todo.Add(t)
			if err != nil {
				log.Println(err)
			}
		}
	}

}
