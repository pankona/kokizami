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
			t, err := todo.Add(v)
			if err != nil {
				log.Println(err)
			}
			log.Println(t)
		}
	}

}
