package command

import (
	"log"

	"github.com/codegangsta/cli"
)

// CmdAdd adds a new todo
func CmdAdd(c *cli.Context) {
	for i, v := range c.Args() {
		log.Println("i =", i)
		log.Println("v =", v)
	}

	log.Println(c.App.Flags)
}
