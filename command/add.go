package command

import (
	"log"

	"github.com/codegangsta/cli"
)

func CmdAdd(c *cli.Context) {
	for i, v := range c.Args() {
		log.Println("i =", i)
		log.Println("v =", v)
	}
}
