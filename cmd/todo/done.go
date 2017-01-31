package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/pankona/todo"
)

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
